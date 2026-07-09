@echo off
echo =======================================================
echo          Welcome to TinyBrain OS + LobeChat UI
echo =======================================================
echo.

REM Start TinyBrain OS Daemon in the background
echo [1/3] Starting TinyBrain OS Engine...
start /B "" ".\tinybrain.exe" daemon --port 8081 > NUL 2>&1

REM Wait a few seconds for the daemon to spin up
timeout /t 3 /nobreak > NUL

REM Start LobeChat Docker Container
echo [2/3] Starting Premium Chat Interface...
docker rm -f tinybrain-chat > NUL 2>&1
docker run -d -p 3210:3210 -e OPENAI_API_KEY=none -e OPENAI_PROXY_URL=http://host.docker.internal:8081/v1 -e OPENAI_MODEL_LIST=-all,+qwen-agent,+sample-alpha --name tinybrain-chat lobehub/lobe-chat > NUL 2>&1

echo [3/3] Opening Browser...
timeout /t 3 /nobreak > NUL
start http://localhost:3210

echo.
echo =======================================================
echo System is ONLINE! You can now chat with your local AI.
echo Close this window to stop the UI. (The Docker container 
echo runs in the background and can be stopped via Docker).
echo =======================================================
pause
