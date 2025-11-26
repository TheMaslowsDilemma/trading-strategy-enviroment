MIGRATIONS_DIR="db/migrate/migrations"
VERSION=$(date +"%Y%m%d%H%M%S")

if [ $# -eq 0 ]; then
  echo "Error: please provide a migration title"
  echo "Usage: $0 <title>"
  exit 1
fi

TITLE="$*"
FILENAME_SAFE=$(echo "$TITLE" | tr '[:upper:]' '[:lower:]' | tr ' ' '_' | sed 's/[^a-z0-9_]//g')

UP_FILE="${MIGRATIONS_DIR}/${VERSION}_${FILENAME_SAFE}.up.sql"
DOWN_FILE="${MIGRATIONS_DIR}/${VERSION}_${FILENAME_SAFE}.down.sql"

cat > "$UP_FILE" <<EOF
-- "$TITLE" Up Migration
-- executed when this migration is applied


EOF

cat > "$DOWN_FILE" <<EOF
-- "$TITLE" Down Migration
-- executed when this migration is rolled back


EOF

echo "New Migration Files:"
echo "  Up:   $UP_FILE"
echo "  Down: $DOWN_FILE"