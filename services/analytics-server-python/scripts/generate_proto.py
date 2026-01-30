from __future__ import annotations

import pathlib
import sys

from grpc_tools import protoc

ROOT = pathlib.Path(__file__).resolve().parents[3]
PROTO_DIR = ROOT / "proto"
OUT_DIR = ROOT / "services" / "analytics-server-python" / "analytics_server" / "proto"


def main() -> int:
    proto_file = PROTO_DIR / "analytics" / "analytics.proto"

    if not proto_file.exists():
        print(f"Proto file not found: {proto_file}")
        return 1

    OUT_DIR.mkdir(parents=True, exist_ok=True)

    result = protoc.main(
        [
            "grpc_tools.protoc",
            f"-I{PROTO_DIR}",
            f"--python_out={OUT_DIR}",
            f"--grpc_python_out={OUT_DIR}",
            str(proto_file),
        ]
    )

    if result != 0:
        print("Failed to generate proto files")
        return result

    _patch_grpc_imports(OUT_DIR / "analytics_pb2_grpc.py")
    print("Proto files generated successfully")
    return 0


def _patch_grpc_imports(path: pathlib.Path) -> None:
    if not path.exists():
        return

    content = path.read_text(encoding="utf-8")
    patched = content.replace("import analytics_pb2 as analytics__pb2", "from . import analytics_pb2 as analytics__pb2")
    if patched != content:
        path.write_text(patched, encoding="utf-8")


if __name__ == "__main__":
    raise SystemExit(main())
