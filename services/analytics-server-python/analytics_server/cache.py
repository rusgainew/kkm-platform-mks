from __future__ import annotations

import json
import logging
from dataclasses import asdict, is_dataclass
from typing import Any, Optional

import redis


class RedisCache:
    def __init__(self, addr: str, password: str, db: int, ttl_seconds: int, logger: logging.Logger):
        self._logger = logger
        self._ttl_seconds = ttl_seconds
        self._enabled = False
        self._client: Optional[redis.Redis] = None

        try:
            if "://" in addr:
                self._client = redis.Redis.from_url(
                    addr,
                    password=password or None,
                    db=db,
                    decode_responses=True,
                )
            else:
                host, _, port = addr.partition(":")
                self._client = redis.Redis(
                    host=host,
                    port=int(port) if port else 6379,
                    password=password or None,
                    db=db,
                    decode_responses=True,
                )
            self._client.ping()
            self._enabled = True
            self._logger.info("Redis connected")
        except Exception as exc:  # noqa: BLE001
            self._logger.warning("Redis unavailable, cache disabled: %s", exc)
            self._enabled = False

    @property
    def enabled(self) -> bool:
        return self._enabled

    def close(self) -> None:
        if self._client is not None:
            try:
                self._client.close()
            except Exception:  # noqa: BLE001
                return

    def get_json(self, key: str) -> Optional[Any]:
        if not self._enabled or self._client is None:
            return None
        try:
            raw = self._client.get(key)
            if not raw:
                return None
            return json.loads(raw)
        except Exception as exc:  # noqa: BLE001
            self._logger.debug("Cache get failed for %s: %s", key, exc)
            return None

    def set_json(self, key: str, value: Any) -> None:
        if not self._enabled or self._client is None:
            return
        try:
            payload = json.dumps(_serialize(value))
            self._client.setex(key, self._ttl_seconds, payload)
        except Exception as exc:  # noqa: BLE001
            self._logger.debug("Cache set failed for %s: %s", key, exc)


def _serialize(value: Any) -> Any:
    if is_dataclass(value):
        return asdict(value)
    if isinstance(value, list):
        return [_serialize(item) for item in value]
    if isinstance(value, dict):
        return {key: _serialize(val) for key, val in value.items()}
    return value
