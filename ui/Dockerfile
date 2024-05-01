FROM node:20-alpine@sha256:7a91aa397f2e2dfbfcdad2e2d72599f374e0b0172be1d86eeb73f1d33f36a4b2 as base

# build step
FROM base as build

RUN mkdir /app
COPY / app/
WORKDIR /app
RUN npm install
RUN npm run build

# final image
FROM base

RUN mkdir -p /app/build
COPY --from=build /app/build /app/package.json /app/entry.js /app/

ENV NODE_ENV=production PORT=8080

ARG GIT_SHA
ENV PUBLIC_VERSION=${GIT_SHA}

CMD ["node", "/app/entry.js"]
