// 自定义校验

// 校验阈值是否在0-100之间
export const checkThreshold = (rule: any, value: any, callback: any) => {
  setTimeout(() => {
    if (!Number.isInteger(parseInt(value))) {
      callback(new Error('请输入数字'))
    } else {
      if (value < 0 || value > 100) {
        callback(new Error('请输入0-100之间的数字'))
      } else {
        callback()
      }
    }
  }, 100)
}

export const checkDuration = (rule: any, value: any, callback: any) => {
  setTimeout(() => {
    if (isNaN(value)) {
      callback(new Error('请输入数字'))
    } else {
      if (value < 0) {
        callback(new Error('请输入大于等于0的数字'))
      } else {
        callback()
      }
    }
  }, 100)
}