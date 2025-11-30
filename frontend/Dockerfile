FROM node:22.15.1 as builder

WORKDIR /app

COPY package*.json ./
RUN yarn install

COPY . .

ARG NEXT_PUBLIC_APP_ENV
ARG NEXT_PUBLIC_BE_API_URL
ENV NEXT_PUBLIC_APP_ENV=$NEXT_PUBLIC_APP_ENV
ENV NEXT_PUBLIC_BE_API_URL=$NEXT_PUBLIC_BE_API_URL

RUN yarn build

FROM nginx:alpine

COPY --from=builder /app/out /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]