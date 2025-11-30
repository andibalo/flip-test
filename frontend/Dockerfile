FROM node:22.15.1 as builder

WORKDIR /app

COPY package*.json ./
RUN yarn install

COPY . .

RUN yarn build

FROM nginx:alpine

COPY --from=builder /app/out /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]