# stop on the first sign of failure
set -e

# function to clean up and stop the emulator
cleanup() {
    echo "Stopping the Spanner emulator..."
    kill -2 $EMULATOR_PID
    docker ps -q --filter "expose=9010-9020" | xargs docker stop
}

# trap any exit signals (EXIT, INT, TERM) and call cleanup function
trap cleanup INT TERM EXIT

# start the Spanner emulator in the background
gcloud emulators spanner start &EMULATOR_PID=$!
echo $EMULATOR_PID

setup_script="./setup-emulator.sh"

# check if the script exists and is executable
if [ -x "$setup_script" ]; then
    echo "Running the Spanner setup script..."
    $setup_script
else
    echo "Error: Setup script is not executable or not found."
    exit 1
fi

echo "Script is running..."
echo "Press Ctrl+C / Control+C to stop."

# Infinite loop to keep the script running
while true; do
    # You can perform your tasks here
    sleep 1  # Sleep for a second to prevent CPU overload
done
