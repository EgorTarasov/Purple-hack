FROM node:18-alpine

WORKDIR /app

COPY . /app

RUN npm install -g pnpm --global-dir=/usr/local/lib/pnpm

RUN pnpm install

RUN pnpm run build

EXPOSE 3000

# RUN npm install  serve
RUN pnpm i  serve
CMD ["pnpm", "run", "serve"]