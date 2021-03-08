FROM python:alpine

ARG PROJECT_DIR=/projects/

# ADD ./requirements.txt $PROJECT_DIR

WORKDIR $PROJECT_DIR
COPY app.py $PROJECT_DIR/

RUN pip install flask

EXPOSE 9000

CMD ["flask","run","--host","0.0.0.0","--port","9000"]
