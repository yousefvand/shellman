# TODO: build.js should work with env to get paths

FROM node:latest

LABEL maintainer="remisa.yousefvand@gmail.com"

RUN mkdir /app
COPY ./build.js /apps
COPY ./package*.json /app
WORKDIR /app
RUN npm i
CMD [ "node", "build.js" ]

# docker build -t remisa/shellman:latest .
# docker push remisa/shellman:latest
# docker run --rm -v "${PWD}":/shellman remisa/shellman:latest
