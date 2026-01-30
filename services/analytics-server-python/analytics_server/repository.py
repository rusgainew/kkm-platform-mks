from __future__ import annotations

import logging
from contextlib import contextmanager
from dataclasses import dataclass
from datetime import datetime
from typing import Iterator, List

import psycopg2
from psycopg2.pool import ThreadedConnectionPool


@dataclass(frozen=True)
class AnalyticsStats:
    total_revenue: float
    total_invoices: int
    average_amount: float
    active_contractors: int
    pending_invoices: int
    approved_invoices: int
    rejected_invoices: int


@dataclass(frozen=True)
class SalesDataPoint:
    date: datetime
    amount: float


@dataclass(frozen=True)
class StatusDistribution:
    status: str
    count: int
    amount: float


@dataclass(frozen=True)
class OperationTypeDistribution:
    operation_type: str
    count: int
    amount: float


@dataclass(frozen=True)
class TopContractorData:
    contractor_id: str
    contractor_name: str
    total_amount: float
    invoice_count: int


@dataclass(frozen=True)
class MonthlyRevenueData:
    month: datetime
    amount: float


class PostgresAnalyticsRepository:
    def __init__(self, dsn: str, max_open_conns: int, max_idle_conns: int, logger: logging.Logger):
        self._logger = logger
        minconn = max(1, max_idle_conns)
        maxconn = max(minconn, max_open_conns)

        self._pool = ThreadedConnectionPool(
            minconn=minconn,
            maxconn=maxconn,
            dsn=dsn,
        )

    def close(self) -> None:
        self._pool.closeall()

    @contextmanager
    def _get_conn(self) -> Iterator[psycopg2.extensions.connection]:
        conn = self._pool.getconn()
        try:
            yield conn
        finally:
            self._pool.putconn(conn)

    def get_stats(self, start_date: datetime, end_date: datetime) -> AnalyticsStats:
        query = """
            SELECT 
                COALESCE(SUM(total_amount), 0) as total_revenue,
                COUNT(*) as total_invoices,
                COALESCE(AVG(total_amount), 0) as average_amount,
                COUNT(DISTINCT contractor_id) as active_contractors,
                COUNT(CASE WHEN status = 'draft' OR status = 'sent' THEN 1 END) as pending_invoices,
                COUNT(CASE WHEN status = 'signed' OR status = 'accepted' THEN 1 END) as approved_invoices,
                COUNT(CASE WHEN status = 'rejected' OR status = 'revoked' THEN 1 END) as rejected_invoices
            FROM invoices
            WHERE created_date >= %s AND created_date <= %s
        """

        with self._get_conn() as conn, conn.cursor() as cur:
            cur.execute(query, (start_date, end_date))
            row = cur.fetchone()
            if row is None:
                return AnalyticsStats(0, 0, 0.0, 0, 0, 0, 0)

            return AnalyticsStats(
                total_revenue=float(row[0] or 0),
                total_invoices=int(row[1] or 0),
                average_amount=float(row[2] or 0),
                active_contractors=int(row[3] or 0),
                pending_invoices=int(row[4] or 0),
                approved_invoices=int(row[5] or 0),
                rejected_invoices=int(row[6] or 0),
            )

    def get_sales_data(self, start_date: datetime, end_date: datetime, granularity: str) -> List[SalesDataPoint]:
        if granularity == "day":
            group = "day"
        elif granularity == "week":
            group = "week"
        elif granularity == "month":
            group = "month"
        else:
            raise ValueError(f"invalid granularity: {granularity}")

        query = f"""
            SELECT 
                DATE_TRUNC('{group}', created_date) as date,
                COALESCE(SUM(total_amount), 0) as amount
            FROM invoices
            WHERE created_date >= %s AND created_date <= %s
            GROUP BY DATE_TRUNC('{group}', created_date)
            ORDER BY date ASC
        """

        with self._get_conn() as conn, conn.cursor() as cur:
            cur.execute(query, (start_date, end_date))
            rows = cur.fetchall()

        return [SalesDataPoint(date=row[0], amount=float(row[1] or 0)) for row in rows]

    def get_status_distribution(self, start_date: datetime, end_date: datetime) -> List[StatusDistribution]:
        query = """
            SELECT 
                status,
                COUNT(*) as count,
                COALESCE(SUM(total_amount), 0) as amount
            FROM invoices
            WHERE created_date >= %s AND created_date <= %s
            GROUP BY status
            ORDER BY count DESC
        """

        with self._get_conn() as conn, conn.cursor() as cur:
            cur.execute(query, (start_date, end_date))
            rows = cur.fetchall()

        return [StatusDistribution(status=row[0], count=int(row[1]), amount=float(row[2] or 0)) for row in rows]

    def get_operation_type_distribution(self, start_date: datetime, end_date: datetime) -> List[OperationTypeDistribution]:
        query = """
            SELECT 
                CASE 
                    WHEN is_resident THEN 'local'
                    ELSE 'import'
                END as operation_type,
                COUNT(*) as count,
                COALESCE(SUM(total_amount), 0) as amount
            FROM invoices
            WHERE created_date >= %s AND created_date <= %s
            GROUP BY is_resident
            ORDER BY count DESC
        """

        with self._get_conn() as conn, conn.cursor() as cur:
            cur.execute(query, (start_date, end_date))
            rows = cur.fetchall()

        return [
            OperationTypeDistribution(operation_type=row[0], count=int(row[1]), amount=float(row[2] or 0))
            for row in rows
        ]

    def get_top_contractors(self, start_date: datetime, end_date: datetime, limit: int) -> List[TopContractorData]:
        query = """
            SELECT 
                contractor_id,
                COALESCE(contractor_id, 'Unknown') as contractor_name,
                COALESCE(SUM(total_amount), 0) as total_amount,
                COUNT(*) as invoice_count
            FROM invoices
            WHERE created_date >= %s AND created_date <= %s
            GROUP BY contractor_id
            ORDER BY total_amount DESC
            LIMIT %s
        """

        with self._get_conn() as conn, conn.cursor() as cur:
            cur.execute(query, (start_date, end_date, limit))
            rows = cur.fetchall()

        return [
            TopContractorData(
                contractor_id=row[0] or "",
                contractor_name=row[1] or "Unknown",
                total_amount=float(row[2] or 0),
                invoice_count=int(row[3] or 0),
            )
            for row in rows
        ]

    def get_monthly_revenue(self, start_date: datetime, end_date: datetime) -> List[MonthlyRevenueData]:
        query = """
            SELECT 
                DATE_TRUNC('month', created_date) as month,
                COALESCE(SUM(total_amount), 0) as amount
            FROM invoices
            WHERE created_date >= %s AND created_date <= %s
            GROUP BY DATE_TRUNC('month', created_date)
            ORDER BY month ASC
        """

        with self._get_conn() as conn, conn.cursor() as cur:
            cur.execute(query, (start_date, end_date))
            rows = cur.fetchall()

        return [MonthlyRevenueData(month=row[0], amount=float(row[1] or 0)) for row in rows]
