# coding: utf-8
from flask import Flask, request, render_template, redirect,jsonify
import os
import socket
import CaboCha
import random
import unittest
import urllib.request
import json
import pykakasi
import MeCab


app = Flask(__name__)
def read_chunk_tokens(tree, chunk):
    """
    チャンクに所属しているトークン列を取得する
    """
    toks = []
    beg = chunk.token_pos
    end = chunk.token_pos + chunk.token_size

    for i in range(beg, end):
        tok = tree.token(i)
        toks.append(tok)

    return toks


def read_subjects_and_predicates(tree):
    """
    主語と述語のリストを読み込む
    """
    subjects = []
    predicates = []

    def backtrack(i):
        """
        チャンクを持っているトークンまで後方検索する
        """
        while i >= 0:
            tok = tree.token(i)
            if tok.chunk:
                return tok
            i -= 1

    for i in range(tree.size()):
        tok = tree.token(i)
        if tok.surface in ('が', 'は', 'を'):
            bef = backtrack(i)
            toks = read_chunk_tokens(tree, bef.chunk)
            subjects.append(toks)
        elif tok.surface in ('た', 'なる', 'ある', 'ない'):
            bef = backtrack(i)
            toks = read_chunk_tokens(tree, bef.chunk)
            predicates.append(toks)

    return subjects, predicates


def merge_surface(toks):
    """
    トークン列の表層形をマージする
    """
    s = ''
    for tok in toks:
        s += tok.surface
    return s


def gen_sentence(subjects, predicates):
    """
    主語のリストと述語のリストからランダムに文章を生成する
    """
    stoks = random.choice(subjects)
    ptoks = random.choice(predicates)
    s1 = merge_surface(stoks)
    s2 = merge_surface(ptoks)
    return s1 + s2

def gen_sentence_hiragana(message,subjects, predicates):
    kakasi = pykakasi.kakasi()
    stoks = random.choice(subjects)
    s1 = merge_surface(stoks)

    
    moji = str.maketrans("ぁぃぅぇぉっゃゅょぢづ", "あいうえおつやゆよじず")
    message=message.translate(moji)
    print(message)
    res = []
    for i in message:
        
        while 1:
            stoks = random.choice(subjects)
            s1 = merge_surface(stoks)

            kakasi.setMode('J', 'H') 
            conv = kakasi.getConverter()
            j_s1 = conv.do(s1)

            kakasi.setMode('K', 'H') 
            conv = kakasi.getConverter()

            j_k_s1 = conv.do(j_s1)
            print(i)
            if(j_k_s1[0]==i):
                break
        ptoks = random.choice(predicates)
        s2 = merge_surface(ptoks)
        kakasi.setMode('J', 'H') 
        conv = kakasi.getConverter()
        j_s2 = conv.do(s2)

        res.append(j_s1 + j_s2)
    return res

def gen_sentence_kanji(message,subjects, predicates):
    kakasi = pykakasi.kakasi()
    stoks = random.choice(subjects)
    s1 = merge_surface(stoks)

    
    moji = str.maketrans("ぁぃぅぇぉっゃゅょぢづ", "あいうえおつやゆよじず")
    message=message.translate(moji)
    print(message)
    res = []
    for i in message:
        
        while 1:
            stoks = random.choice(subjects)
            s1 = merge_surface(stoks)

            kakasi.setMode('J', 'H') 
            conv = kakasi.getConverter()
            j_s1 = conv.do(s1)

            kakasi.setMode('K', 'H') 
            conv = kakasi.getConverter()

            j_k_s1 = conv.do(j_s1)
            print(i)
            if(j_k_s1[0]==i):
                break
        ptoks = random.choice(predicates)
        s2 = merge_surface(ptoks)
        kakasi.setMode('J', 'H') 
        conv = kakasi.getConverter()
        j_s2 = conv.do(s2)

        res.append(s1 + s2)
    return res

@app.route("/message",methods=['POST'])
def hello():
       # jsonレスポンス返却
       # jsonifyにdict型オブジェクトを設定するとjsonデータのレスポンスが生成される
    f = open('./text.txt', 'r')
    sentence = f.read()
    f.close()
    cp = CaboCha.Parser()
    tree = cp.parse(sentence)
    g_subjects, g_predicates = read_subjects_and_predicates(tree)
    message = request.json['message']
    result = gen_sentence_hiragana(message,g_subjects, g_predicates)
    #後でスコアを出す関数に変える
    score = 114514 
    return jsonify({'messages': result,"score":score})

@app.route("/message2",methods=['POST'])
def hello2():
       # jsonレスポンス返却
       # jsonifyにdict型オブジェクトを設定するとjsonデータのレスポンスが生成される
    f = open('./text.txt', 'r')
    sentence = f.read()
    f.close()
    cp = CaboCha.Parser()
    tree = cp.parse(sentence)
    g_subjects, g_predicates = read_subjects_and_predicates(tree)
    message = request.json['message']
    result = gen_sentence_kanji(message,g_subjects, g_predicates)
    #後でスコアを出す関数に変える
    score = 114514 
    return jsonify({'messages': result,"score":score})

@app.route("/test",methods=['GET'])
def test():
    return jsonify({'messages': ["あっぷる","いか","しかばね","てがみ","るびー"]})

    #NLP（自然言語処理）
    
    """
    青空文庫から文字を取得
    """
    # url = 'http://www.aozorahack.net/api/v0.1/books/773/content?format=txt'

    # try:
    #   with urllib.request.urlopen(url) as response:
    #         print(response.read().decode(encoding='shift_jis', errors='replace'))
    #         body = json.loads(response.read().decode(encoding='shift_jis', errors='replace'))
    #         headers = response.getheaders()
    #         status = response.getcode()

    #         print(headers)
    #         print(status)
    #         print(body)

    # except urllib.error.URLError as e:
    #     print(e)#e.reason        
    
    

    # 漢字が入っている文をひらがなだけの文に変換
    #conv = kakasi.getConverter()

    # print("----------主語----------")
    # for s in subjects:
    #     print(conv.do(merge_surface(s)))
    # print("-----------------------")

    # print("----------述語----------")
    # for p in predicates:
    #     print(conv.do(merge_surface(p)))
    # print("-----------------------")

    
    # for _ in range(20):
    #     result = gen_sentence(subjects, predicates)
    #     print(result)

class Test(unittest.TestCase):
    def eq(self, a, b, c):
        cp = CaboCha.Parser()
        tree = cp.parse(a)
        subjects, predicates = read_subjects_and_predicates(tree)
        for i in range(len(subjects)):
            s = merge_surface(subjects[i])
            self.assertEqual(s, b[i])
        for i in range(len(predicates)):
            s = merge_surface(predicates[i])
            self.assertEqual(s, c[i])

    def test_method(self):
        self.eq('猫は犬である', ['猫は'], ['犬である'])
        self.eq('可愛い猫は怖い犬である', ['猫は'], ['犬である'])
        self.eq('猫が話した', ['猫が'], ['話した'])
        self.eq('猫が大いに話した', ['猫が'], ['話した'])
        self.eq('猫を可愛がらない', ['猫を'], ['可愛がらない'])
        self.eq('猫が鳥になる', ['猫が'], ['なる'])
        self.eq('猫は眠らない', ['猫は'], ['眠らない'])



if __name__ == "__main__":

    app.run(host='0.0.0.0', port=9000)
