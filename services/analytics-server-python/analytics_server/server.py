from __future__ import annotations

import signal
from concurrent import futures

import grpc
from grpc_health.v1 import health, health_pb2, health_pb2_grpc

from .cache import RedisCache
from .config import build_dsn, load_config
from .grpc_handlers import AnalyticsServicer
from .logging_setup import configure_logging
from .repository import PostgresAnalyticsRepository
from .service import AnalyticsService

try:
    from analytics_server.proto import analytics_pb2_grpc
except ImportError as exc:  # pragma: no cover
    raise RuntimeError(
        "Proto files not generated. Run 'python scripts/generate_proto.py' first."
    ) from exc


def serve() -> None:
    config = load_config()
    logger = configure_logging(config.log_level)

    logger.info("Starting analytics-server", extra={"grpc_port": config.server.grpc_port})

    dsn = build_dsn(config.database)
    repo = PostgresAnalyticsRepository(
        dsn=dsn,
        max_open_conns=config.database.max_open_conns,
        max_idle_conns=config.database.max_idle_conns,
        logger=logger,
    )

    cache = RedisCache(
        addr=config.redis.addr,
        password=config.redis.password,
        db=config.redis.db,
        ttl_seconds=config.redis.ttl_seconds,
        logger=logger,
    )

    service = AnalyticsService(repo=repo, cache=cache)

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    analytics_pb2_grpc.add_AnalyticsServiceServicer_to_server(
        AnalyticsServicer(service=service, logger=logger),
        server,
    )

    health_servicer = health.HealthServicer()
    health_pb2_grpc.add_HealthServicer_to_server(health_servicer, server)
    health_servicer.set("", health_pb2.HealthCheckResponse.SERVING)

    server.add_insecure_port(f"[::]:{config.server.grpc_port}")
    server.start()

    def _graceful_shutdown(signum, frame):  # noqa: ANN001
        logger.info("Shutting down analytics-server")
        server.stop(3)
        cache.close()
        repo.close()

    signal.signal(signal.SIGINT, _graceful_shutdown)
    signal.signal(signal.SIGTERM, _graceful_shutdown)

    server.wait_for_termination()


if __name__ == "__main__":
    serve()
