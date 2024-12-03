/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-prometheus licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Jul 26 16:42:38 2023 +0800
 */
module.exports = {
  devServer: {
    headers: {
      'Access-Control-Allow-Origin': '*',
    },
    proxy: {
      '/': {
        target: 'http://localhost:8090',
        ws: false,
        changeOrigin: true,
      }
    },
    port: 8082,
    host: 'localhost',
  },

  // 静态资源
  publicPath: "./prometheus",
  outputDir: "dist",
  assetsDir: "static",
  indexPath: "index.html",
  filenameHashing: true,

  chainWebpack: (config) => {
    config.resolve.symlinks(true);
  },
  configureWebpack: {
    output: {
      //资源打包路径
      library: 'vueApp',
      libraryTarget: 'umd'
    }
  }
}
