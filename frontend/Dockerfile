FROM node:8.6
EXPOSE 8080

WORKDIR /app
ADD package.json package.json
ADD yarn.lock yarn.lock
RUN yarn
ADD .babelrc .babelrc

ADD index.html index.html
ADD webpack.config.js webpack.config.js

ADD src src

CMD ["npm", "run", "dev"]
