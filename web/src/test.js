
    const item= {
      metric: {
        fstype: 'ext4'
      },
      value: [23, 33]
    }
  
  
  // 要访问的路径
  const path1 = 'metric.fstype'; // 访问 fstype
  const path2 = 'value[1]';       // 访问 value 数组中的第一个元素
  
  // 使用 split 和 reduce 访问嵌套属性的函数
  const getValueFromPath = (obj, path) => {
    return path.split(/\.|\[|\]/) // 使用正则表达式分割字符串
               .filter(Boolean)   // 过滤掉空字符串
               .reduce((acc, part) => acc && acc[part], obj); // 逐层访问
  };
  
  // 访问 fstype
  const fstypeValue = getValueFromPath(item, path1);
  console.log(fstypeValue); // 输出: 'ext4'
  
  // 访问 value 数组中的第一个元素
  const valueElement = getValueFromPath(item, path2);
  console.log(valueElement); // 输出: 23
  