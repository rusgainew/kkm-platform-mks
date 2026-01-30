from __future__ import annotations

import os
from dataclasses import dataclass
from datetime import timedelta
from typing import Optional

from dotenv import load_dotenv

load_dotenv()


def _get_env(key: str, default: str) -> str:
    value = os.getenv(key)
    return value if value not in (None, "") else default


def _get_env_int(key: str, default: int) -> int:
    value = os.getenv(key)
    if not value:
        return default
    try:
        return int(value)
    except ValueError:
        return default


def _parse_duration(value: str, default: timedelta) -> timedelta:
    if not value:
        return default

    value = value.strip().lower()
    try:
        if value.endswith("ms"):
            return timedelta(milliseconds=int(value[:-2]))
        if value.endswith("s"):
            return timedelta(seconds=int(value[:-1]))
        if value.endswith("m"):
            return timedelta(minutes=int(value[:-1]))
        if value.endswith("h"):
            return timedelta(hours=int(value[:-1]))
        return timedelta(seconds=int(value))
    except ValueError:
        return default


@dataclass(frozen=True)
class ServerConfig:
    grpc_port: int
    metrics_port: int


@dataclass(frozen=True)
class DatabaseConfig:
    host: str
    port: int
    user: str
    password: str
    dbname: str
    sslmode: str
    max_open_conns: int
    max_idle_conns: int
    conn_max_lifetime: timedelta


@dataclass(frozen=True)
class RedisConfig:
    addr: str
    password: str
    db: int
    ttl_seconds: int


@dataclass(frozen=True)
class AppConfig:
    server: ServerConfig
    database: DatabaseConfig
    redis: RedisConfig
    log_level: str


def load_config() -> AppConfig:
    grpc_port = _get_env_int("GRPC_PORT", 50070)
    metrics_port = _get_env_int("METRICS_PORT", 9115)

    db_conn_max_lifetime = _parse_duration(
        _get_env("DB_CONN_MAX_LIFETIME", "300s"),
        timedelta(minutes=5),
    )

    config = AppConfig(
        server=ServerConfig(
            grpc_port=grpc_port,
            metrics_port=metrics_port,
        ),
        database=DatabaseConfig(
            host=_get_env("DB_HOST", "localhost"),
            port=_get_env_int("DB_PORT", 5432),
            user=_get_env("DB_USER", "analytics_svc"),
            password=_get_env("DB_PASSWORD", "analytics_pass_2026"),
            dbname=_get_env("DB_NAME", "analytics_db"),
            sslmode=_get_env("DB_SSLMODE", "disable"),
            max_open_conns=_get_env_int("DB_MAX_OPEN_CONNS", 25),
            max_idle_conns=_get_env_int("DB_MAX_IDLE_CONNS", 5),
            conn_max_lifetime=db_conn_max_lifetime,
        ),
        redis=RedisConfig(
            addr=_get_env("REDIS_ADDR", "localhost:6379"),
            password=_get_env("REDIS_PASSWORD", ""),
            db=_get_env_int("REDIS_DB", 0),
            ttl_seconds=_get_env_int("REDIS_TTL", 300),
        ),
        log_level=_get_env("LOG_LEVEL", "INFO").upper(),
    )

    if not config.database.host:
        raise ValueError("DB_HOST is required")
    if not config.database.dbname:
        raise ValueError("DB_NAME is required")

    return config


def build_dsn(cfg: DatabaseConfig) -> str:
    return (
        "host={host} port={port} user={user} password={password} "
        "dbname={dbname} sslmode={sslmode}"
    ).format(
        host=cfg.host,
        port=cfg.port,
        user=cfg.user,
        password=cfg.password,
        dbname=cfg.dbname,
        sslmode=cfg.sslmode,
    )
