import axios from 'axios'

const GEO_DATAV_URL = 'https://geo.datav.aliyun.com/areas_v3/bound/100000_full.json'

export const getChinaGeoJSON = async () => {
  const response = await axios.get(GEO_DATAV_URL)
  return response.data
}
