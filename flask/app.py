# coding: utf-8
from flask import Flask
from flask import jsonify
import os
import socket

app = Flask(__name__)

@app.route("/message",methods=['POST'])
def hello():
    # jsonレスポンス返却
    # jsonifyにdict型オブジェクトを設定するとjsonデータのレスポンスが生成される
    return jsonify({'messages': ["あっぷる","いか","しかばね","てがみ","るびー"],"score":114514})

@app.route("/test",methods=['GET'])
def test():
    return jsonify({'messages': ["あっぷる","いか","しかばね","てがみ","るびー"]})

if __name__ == "__main__":
  app.run(host='0.0.0.0', port=9000)
