from __future__ import annotations

import logging
from datetime import datetime

import grpc
from google.protobuf.timestamp_pb2 import Timestamp

from .service import AnalyticsService

try:
    from analytics_server.proto import analytics_pb2, analytics_pb2_grpc
except ImportError as exc:  # pragma: no cover
    raise RuntimeError(
        "Proto files not generated. Run 'python scripts/generate_proto.py' first."
    ) from exc


class AnalyticsServicer(analytics_pb2_grpc.AnalyticsServiceServicer):
    def __init__(self, service: AnalyticsService, logger: logging.Logger):
        self._service = service
        self._logger = logger

    def GetDashboardStats(self, request, context):  # noqa: N802
        start_date = _ts_to_dt(request.start_date)
        end_date = _ts_to_dt(request.end_date)

        stats = self._service.get_stats(start_date, end_date)
        return analytics_pb2.StatsResponse(
            total_invoices=stats.total_invoices,
            total_revenue=stats.total_revenue,
            average_amount=stats.average_amount,
            unique_contractors=stats.active_contractors,
            pending_count=stats.pending_invoices,
            approved_count=stats.approved_invoices,
            rejected_count=stats.rejected_invoices,
        )

    def GetSalesChart(self, request, context):  # noqa: N802
        start_date = _ts_to_dt(request.start_date)
        end_date = _ts_to_dt(request.end_date)

        data = self._service.get_sales_data(start_date, end_date, request.granularity)
        points = [
            analytics_pb2.SalesDataPoint(
                date=item.date.strftime("%Y-%m-%d"),
                count=1,
                amount=item.amount,
            )
            for item in data
        ]
        return analytics_pb2.SalesChartResponse(data_points=points)

    def GetStatusDistribution(self, request, context):  # noqa: N802
        start_date = _ts_to_dt(request.start_date)
        end_date = _ts_to_dt(request.end_date)

        data = self._service.get_status_distribution(start_date, end_date)
        total_amount = sum(item.amount for item in data) or 0
        items = []
        for item in data:
            percentage = (item.amount / total_amount * 100) if total_amount > 0 else 0
            items.append(
                analytics_pb2.StatusDistributionItem(
                    status=item.status,
                    count=item.count,
                    percentage=percentage,
                )
            )

        return analytics_pb2.StatusResponse(items=items)

    def GetOperationTypeDistribution(self, request, context):  # noqa: N802
        start_date = _ts_to_dt(request.start_date)
        end_date = _ts_to_dt(request.end_date)

        data = self._service.get_operation_type_distribution(start_date, end_date)
        total_amount = sum(item.amount for item in data) or 0
        items = []
        for item in data:
            percentage = (item.amount / total_amount * 100) if total_amount > 0 else 0
            items.append(
                analytics_pb2.OperationTypeItem(
                    operation_type=item.operation_type,
                    count=item.count,
                    amount=item.amount,
                    percentage=percentage,
                )
            )

        return analytics_pb2.OperationTypeResponse(items=items)

    def GetTopContractors(self, request, context):  # noqa: N802
        start_date = _ts_to_dt(request.start_date)
        end_date = _ts_to_dt(request.end_date)
        limit = request.limit or 10

        data = self._service.get_top_contractors(start_date, end_date, limit)
        contractors = [
            analytics_pb2.TopContractorItem(
                contractor_id=item.contractor_id,
                contractor_name=item.contractor_name,
                invoice_count=item.invoice_count,
                total_amount=item.total_amount,
            )
            for item in data
        ]

        return analytics_pb2.TopContractorsResponse(contractors=contractors)

    def GetMonthlyRevenue(self, request, context):  # noqa: N802
        start_date = _ts_to_dt(request.start_date)
        end_date = _ts_to_dt(request.end_date)

        data = self._service.get_monthly_revenue(start_date, end_date)
        months = [
            analytics_pb2.MonthlyRevenueItem(
                month=item.month.strftime("%Y-%m"),
                revenue=item.amount,
                invoice_count=1,
            )
            for item in data
        ]

        return analytics_pb2.MonthlyRevenueResponse(months=months)


def _ts_to_dt(ts: Timestamp) -> datetime:
    return ts.ToDatetime()
