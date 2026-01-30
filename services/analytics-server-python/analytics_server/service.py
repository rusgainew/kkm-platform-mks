from __future__ import annotations

from dataclasses import asdict
from datetime import datetime
from typing import List

from .cache import RedisCache
from .repository import (
    AnalyticsStats,
    MonthlyRevenueData,
    OperationTypeDistribution,
    SalesDataPoint,
    StatusDistribution,
    TopContractorData,
    PostgresAnalyticsRepository,
)


class AnalyticsService:
    def __init__(self, repo: PostgresAnalyticsRepository, cache: RedisCache | None):
        self._repo = repo
        self._cache = cache

    def get_stats(self, start_date: datetime, end_date: datetime) -> AnalyticsStats:
        cache_key = self._cache_key("stats", start_date, end_date)
        cached = self._get_cached(cache_key)
        if cached is not None:
            return AnalyticsStats(**cached)

        stats = self._repo.get_stats(start_date, end_date)
        self._set_cached(cache_key, asdict(stats))
        return stats

    def get_sales_data(self, start_date: datetime, end_date: datetime, granularity: str) -> List[SalesDataPoint]:
        cache_key = self._cache_key("sales", start_date, end_date, granularity)
        cached = self._get_cached(cache_key)
        if cached is not None:
            return [SalesDataPoint(date=_parse_dt(item["date"]), amount=item["amount"]) for item in cached]

        data = self._repo.get_sales_data(start_date, end_date, granularity)
        payload = [
            {"date": item.date.isoformat(), "amount": item.amount}
            for item in data
        ]
        self._set_cached(cache_key, payload)
        return data

    def get_status_distribution(self, start_date: datetime, end_date: datetime) -> List[StatusDistribution]:
        cache_key = self._cache_key("status_dist", start_date, end_date)
        cached = self._get_cached(cache_key)
        if cached is not None:
            return [StatusDistribution(**item) for item in cached]

        data = self._repo.get_status_distribution(start_date, end_date)
        self._set_cached(cache_key, [asdict(item) for item in data])
        return data

    def get_operation_type_distribution(
        self, start_date: datetime, end_date: datetime
    ) -> List[OperationTypeDistribution]:
        cache_key = self._cache_key("operation_dist", start_date, end_date)
        cached = self._get_cached(cache_key)
        if cached is not None:
            return [OperationTypeDistribution(**item) for item in cached]

        data = self._repo.get_operation_type_distribution(start_date, end_date)
        self._set_cached(cache_key, [asdict(item) for item in data])
        return data

    def get_top_contractors(
        self, start_date: datetime, end_date: datetime, limit: int
    ) -> List[TopContractorData]:
        cache_key = self._cache_key("top_contractors", start_date, end_date, str(limit))
        cached = self._get_cached(cache_key)
        if cached is not None:
            return [TopContractorData(**item) for item in cached]

        data = self._repo.get_top_contractors(start_date, end_date, limit)
        self._set_cached(cache_key, [asdict(item) for item in data])
        return data

    def get_monthly_revenue(self, start_date: datetime, end_date: datetime) -> List[MonthlyRevenueData]:
        cache_key = self._cache_key("monthly_revenue", start_date, end_date)
        cached = self._get_cached(cache_key)
        if cached is not None:
            return [MonthlyRevenueData(month=_parse_dt(item["month"]), amount=item["amount"]) for item in cached]

        data = self._repo.get_monthly_revenue(start_date, end_date)
        payload = [
            {"month": item.month.isoformat(), "amount": item.amount}
            for item in data
        ]
        self._set_cached(cache_key, payload)
        return data

    def _cache_key(self, prefix: str, start_date: datetime, end_date: datetime, *parts: str) -> str:
        start_str = start_date.strftime("%Y-%m-%d")
        end_str = end_date.strftime("%Y-%m-%d")
        extra = ":".join(parts)
        if extra:
            return f"analytics:{prefix}:{start_str}:{end_str}:{extra}"
        return f"analytics:{prefix}:{start_str}:{end_str}"

    def _get_cached(self, key: str):
        if not self._cache or not self._cache.enabled:
            return None
        return self._cache.get_json(key)

    def _set_cached(self, key: str, value):
        if not self._cache or not self._cache.enabled:
            return
        self._cache.set_json(key, value)


def _parse_dt(value: str) -> datetime:
    return datetime.fromisoformat(value)
