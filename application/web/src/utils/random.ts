// 随机生成中文姓名
const lastNames = ['张', '王', '李', '赵', '刘', '陈', '杨', '黄', '周', '吴'];
const firstNames = ['伟', '芳', '娜', '秀英', '敏', '静', '丽', '强', '磊', '洋', '艳', '勇', '军', '杰', '娟', '涛', '明', '超', '秀兰', '霞'];

export const generateRandomName = () => {
  const lastName = lastNames[Math.floor(Math.random() * lastNames.length)];
  const firstName = firstNames[Math.floor(Math.random() * firstNames.length)];
  return lastName + firstName;
};

// 随机生成地址
const cities = ['北京市', '上海市', '广州市', '深圳市', '杭州市', '南京市', '成都市', '武汉市'];
const districts = ['东城区', '西城区', '朝阳区', '海淀区', '丰台区', '昌平区'];
const streets = ['长安街', '建国路', '复兴路', '中关村大街', '金融街', '望京街'];
const communities = ['阳光小区', '和平花园', '翠湖园', '金色家园', '龙湖花园', '碧桂园'];

export const generateRandomAddress = () => {
  const city = cities[Math.floor(Math.random() * cities.length)];
  const district = districts[Math.floor(Math.random() * districts.length)];
  const street = streets[Math.floor(Math.random() * streets.length)];
  const community = communities[Math.floor(Math.random() * communities.length)];
  const building = Math.floor(Math.random() * 20 + 1);
  const unit = Math.floor(Math.random() * 6 + 1);
  const room = Math.floor(Math.random() * 2000 + 101);

  return `${city}${district}${street}${community}${building}号楼${unit}单元${room}室`;
};

// 随机生成面积（50-300平方米）
export const generateRandomArea = () => {
  return Number((Math.random() * (300 - 50) + 50).toFixed(2));
};

// 随机生成价格（50-1000万）
export const generateRandomPrice = () => {
  return Number((Math.random() * (10000000 - 500000) + 500000).toFixed(2));
}; 