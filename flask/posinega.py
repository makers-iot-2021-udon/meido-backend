from bs4 import BeautifulSoup
# ファイルを読み込む
with open("syosetu.html", "rt", encoding="sjis") as f:
  html = f.read()
  # HTMLをパースする
  soup = BeautifulSoup(html, 'html.parser')
  # ルビを削除
  soup.find("rp").extract()
  soup.find("rt").extract()
  # テキストだけを取り出す
  text = soup.get_text()
  print(text)
  # 保存
  with open("syosetu.txt", "wt", encoding="utf-8") as w:
    w.write(text)