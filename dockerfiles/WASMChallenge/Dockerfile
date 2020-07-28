FROM ubuntu AS build_emcc
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get -qq update

RUN apt-get install -qq -y \
    python3.8 cmake default-jre git-core

RUN ln -s /usr/bin/python3.8 /usr/bin/python

RUN git clone https://github.com/emscripten-core/emsdk.git

RUN cd /emsdk \
 && git pull \
 && ./emsdk install latest \
 && ./emsdk activate latest \
 && bash -c '. ./emsdk_env.sh'

FROM build_emcc AS build_flag

#FLAG as evironment variable
ENV FLAG BSidesTLV2020{W1sdom_i5_only_pOs5ess3d_by_th3_l34rned}

COPY ctf_data /root/webapp

WORKDIR /root/webapp

#Compile with new flag
RUN /emsdk/upstream/emscripten/emcc -DFLAG_BYTES=$(python ./encode_flag.py) \
    -s -Wl,--no-entry snake.c -o snake.wasm --js-library snake.js

RUN mkdir -p html && cp Phone.png index.html snake.css snake.js snake.wasm html

FROM nginx:alpine 

COPY --from=build_flag /root/webapp/html /usr/share/nginx/html

RUN sed -i '/\}/i application/wasm wasm;' /etc/nginx/mime.types

CMD ["nginx", "-g", "daemon off;"]
