FROM alpine:3.12 as build

COPY . /app
WORKDIR /app

# Install packages 
RUN apk update && apk add gcc g++

# Install 
RUN gcc -c -std=c99 zlib/*.c 
RUN g++ -c -std=c++14 snappaste.cc 
RUN g++ -o snappaste *.o -pthread


FROM node:12.18-alpine as core

WORKDIR /app

# Install packages 
RUN npm i -g pm2

# Set non root user
RUN addgroup -S bsidestlv && adduser -S bsidestlv -G bsidestlv

# Copy core
COPY --from=build /app/paste/ /app/paste/
COPY --from=build /app/www/ /app/www/
COPY --from=build /app/flag /app/snappaste /app/

RUN chown -R bsidestlv:bsidestlv /app/
USER bsidestlv

EXPOSE 8080
CMD [ "pm2", "start", "--no-daemon", "/app/snappaste" ]
