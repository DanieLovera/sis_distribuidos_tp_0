FROM python:3.9.7-slim
COPY server /
RUN python3 -m unittest discover -s test -p "test_*.py" -v
ENTRYPOINT ["/bin/sh"]
