FROM python:3.7-alpine

RUN apk update
RUN apk --no-cache add git gcc curl bash file sudo openssh build-base openssl swig

# Install CRF++
ENV version 0.58
ENV url https://drive.google.com/uc?export=download&id=0B4y35FiV1wh7QVR6VXJ5dWExSTQ
RUN cd /tmp \
    && curl -L -o CRF++-${version}.tar.gz ${url} \
    && tar -zxf CRF++-${version}.tar.gz \
    && cd CRF++-${version} \
    && ./configure \
    && make \
    && make install \
    && rm -rf /tmp/*

# Install MeCab
ENV version 0.996
ENV url https://drive.google.com/uc?export=download&id=0B4y35FiV1wh7cENtOXlicTFaRUE
RUN cd /tmp \
    && curl -L -o mecab-${version}.tar.gz ${url} \
    && tar -zxf mecab-${version}.tar.gz \
    && cd mecab-${version} \
    && ./configure --enable-utf8-only --with-charset=utf8 \
    && make \
    && make install \
    && rm -rf /tmp/*

# Install IPA dic
ENV version 2.7.0-20070801
ENV url https://drive.google.com/uc?export=download&id=0B4y35FiV1wh7MWVlSDBCSXZMTXM
RUN cd /tmp \
    && curl -SL -o mecab-ipadic-${version}.tar.gz ${url} \
    && tar -zxf mecab-ipadic-${version}.tar.gz \
    && cd mecab-ipadic-${version} \
    && ./configure --with-charset=utf8 \
    && make \
    && make install \
    && rm -rf /tmp/*

# Install mecab-ipadic-NEologd
RUN cd /tmp \
    && git clone --depth 1 https://github.com/neologd/mecab-ipadic-neologd.git \
    && mecab-ipadic-neologd/bin/install-mecab-ipadic-neologd -n -y \
    && rm -rf /tmp/*

# Install CaboCha
ENV version 0.69
ENV url https://drive.google.com/uc?export=download&id=0B4y35FiV1wh7SDd1Q1dUQkZQaUU
ENV url2 https://drive.google.com/uc?confirm={}&export=download&id=0B4y35FiV1wh7SDd1Q1dUQkZQaUU
ENV CPPFLAGS -I/usr/local/include
RUN curl -c cabocha-${version}.tar.bz2 -s -L ${url} \
    | grep confirm \
    | sed -e "s/^.*confirm=\(.*\)&amp;id=.*$/\1/" \
    | xargs -I{} curl -b  cabocha-${version}.tar.bz2 -L -o cabocha-${version}.tar.bz2 ${url2} \
    && tar -jxf cabocha-${version}.tar.bz2 \
    && cd cabocha-${version} \
    && ./configure --prefix=/usr/local --with-charset=utf8 \
    && make \
    && make install \
    && rm -rf /tmp/*

RUN pip3 install mecab-python3
RUN pip3 install git+https://github.com/kenkov/cabocha
RUN pip3 install cabocha-python
#追加分
RUN pip3 install pykakasi