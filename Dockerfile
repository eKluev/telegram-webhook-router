FROM python:3.9
WORKDIR /usr/src/app/
COPY . /usr/src/app/
RUN pip install -r requirements.txt
EXPOSE 80
CMD ["gunicorn", "--bind", "0.0.0.0:80", "app:app"]