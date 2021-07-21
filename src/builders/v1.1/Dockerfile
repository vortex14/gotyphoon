FROM {{.TYPHOON_IMAGE}}
WORKDIR /app
COPY ["donor.py", "fetcher.py", "processor.py", "result_transporter.py", "scheduler.py", "config.kube-stage.yaml", "/app/"]
COPY ["./project", "/app/project"]
CMD ["python", "donor.py", "--config=config.kube-stage.yaml", "--level=DEBUG"]