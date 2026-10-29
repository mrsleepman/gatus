/**
 * Generates a human-readable relative time string (e.g., "2 hours ago")
 * @param {string|Date} timestamp - The timestamp to convert
 * @returns {string} Relative time string
 */
export const generatePrettyTimeAgo = (timestamp) => {
  let differenceInMs = new Date().getTime() - new Date(timestamp).getTime();
  if (differenceInMs < 500) {
    return "刚刚";
  }
  if (differenceInMs > 3 * 86400000) {
    let days = (differenceInMs / 86400000).toFixed(0);
    return days + "天前";
  }
  if (differenceInMs > 3600000) {
    let hours = (differenceInMs / 3600000).toFixed(0);
    return hours + "小时前";
  }
  if (differenceInMs > 60000) {
    let minutes = (differenceInMs / 60000).toFixed(0);
    return minutes + "分钟前";
  }
  let seconds = (differenceInMs / 1000).toFixed(0);
  return seconds + "秒前";
}

/**
 * Generates a pretty time difference string between two timestamps
 * @param {string|Date} start - Start timestamp
 * @param {string|Date} end - End timestamp
 * @returns {string} Time difference string
 */
export const generatePrettyTimeDifference = (start, end) => {
  const ms = new Date(start) - new Date(end)
  const seconds = Math.floor(ms / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)

  if (hours > 0) {
    const remainingMinutes = minutes % 60
    const hoursText = hours + '小时'
    if (remainingMinutes > 0) {
      return hoursText + remainingMinutes + '分钟'
    }
    return hoursText
  } else if (minutes > 0) {
    const remainingSeconds = seconds % 60
    const minutesText = minutes + '分钟'
    if (remainingSeconds > 0) {
      return minutesText + remainingSeconds + '秒'
    }
    return minutesText
  } else {
    return seconds + '秒'
  }
}

/**
 * Formats a timestamp into YYYY-MM-DD HH:mm:ss format
 * @param {string|Date} timestamp - The timestamp to format
 * @returns {string} Formatted timestamp
 */
export const prettifyTimestamp = (timestamp) => {
  let date = new Date(timestamp);
  let YYYY = date.getFullYear();
  let MM = ((date.getMonth() + 1) < 10 ? "0" : "") + "" + (date.getMonth() + 1);
  let DD = ((date.getDate()) < 10 ? "0" : "") + "" + (date.getDate());
  let hh = ((date.getHours()) < 10 ? "0" : "") + "" + (date.getHours());
  let mm = ((date.getMinutes()) < 10 ? "0" : "") + "" + (date.getMinutes());
  let ss = ((date.getSeconds()) < 10 ? "0" : "") + "" + (date.getSeconds());
  return YYYY + "-" + MM + "-" + DD + " " + hh + ":" + mm + ":" + ss;
}