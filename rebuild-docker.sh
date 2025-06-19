
set -e

echo "Removing previous containers and volumes..."
docker compose down -v

echo "Rebuilding containers..."
docker compose up --build

echo "Docker services successfully restarted!"