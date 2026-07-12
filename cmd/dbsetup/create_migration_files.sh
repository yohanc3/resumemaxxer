migrate create -ext sql -dir migration_files -seq $1

if [ $? -eq 0 ]; then
    echo "Successfully created both files."
else
    echo "Something went wrong. Exit code $"
fi
