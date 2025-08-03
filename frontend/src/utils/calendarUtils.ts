/**
 * 生成ICS格式的日历事件
 */

// ICS文件头
const ICS_HEADER = [
  'BEGIN:VCALENDAR',
  'VERSION:2.0',
  'PRODID:-//todoIng//Calendar//EN',
  'CALSCALE:GREGORIAN'
].join('\r\n');

// ICS文件尾
const ICS_FOOTER = 'END:VCALENDAR';

/**
 * 格式化日期为ICS格式
 * @param date 日期字符串
 * @returns ICS格式的日期字符串
 */
const formatICalendarDate = (date: string): string => {
  const d = new Date(date);
  const year = d.getFullYear();
  const month = (d.getMonth() + 1).toString().padStart(2, '0');
  const day = d.getDate().toString().padStart(2, '0');
  const hours = d.getHours().toString().padStart(2, '0');
  const minutes = d.getMinutes().toString().padStart(2, '0');
  const seconds = d.getSeconds().toString().padStart(2, '0');
  return `${year}${month}${day}T${hours}${minutes}${seconds}`;
};

/**
 * 转义ICS文本中的特殊字符
 * @param text 需要转义的文本
 * @returns 转义后的文本
 */
const escapeICalendarText = (text: string): string => {
  return text
    .replace(/\\/g, '\\\\')
    .replace(/;/g, '\\;')
    .replace(/,/g, '\\,')
    .replace(/\n/g, '\\n');
};

/**
 * 生成单个任务的ICS事件
 * @param task 任务对象
 * @returns ICS事件字符串
 */
export const generateTaskEvent = (task: any): string => {
  const lines = ['BEGIN:VEVENT'];
  
  // 事件UID
  lines.push(`UID:${task._id}@todoing`);
  
  // 创建时间
  lines.push(`DTSTAMP:${formatICalendarDate(new Date().toISOString())}`);
  
  // 任务标题
  lines.push(`SUMMARY:${escapeICalendarText(task.title)}`);
  
  // 任务描述
  if (task.description) {
    lines.push(`DESCRIPTION:${escapeICalendarText(task.description)}`);
  }
  
  // 截止日期作为事件日期
  if (task.deadline) {
    lines.push(`DTSTART:${formatICalendarDate(task.deadline)}`);
    lines.push(`DTEND:${formatICalendarDate(new Date(new Date(task.deadline).getTime() + 3600000).toISOString())}`);
  } else if (task.scheduledDate) {
    // 如果没有截止日期，则使用计划日期
    lines.push(`DTSTART:${formatICalendarDate(task.scheduledDate)}`);
    lines.push(`DTEND:${formatICalendarDate(new Date(new Date(task.scheduledDate).getTime() + 3600000).toISOString())}`);
  }
  
  // 任务状态
  lines.push(`STATUS:${task.status === 'Done' ? 'COMPLETED' : task.status === 'In Progress' ? 'CONFIRMED' : 'TENTATIVE'}`);
  
  // 优先级
  if (task.priority === 'High') {
    lines.push('PRIORITY:1');
  } else if (task.priority === 'Medium') {
    lines.push('PRIORITY:5');
  } else if (task.priority === 'Low') {
    lines.push('PRIORITY:9');
  }
  
  lines.push('END:VEVENT');
  return lines.join('\r\n');
};

/**
 * 生成多个任务的ICS日历文件内容
 * @param tasks 任务数组
 * @returns ICS日历文件内容
 */
export const generateCalendarICS = (tasks: any[]): string => {
  const events = tasks
    .filter(task => task.deadline || task.scheduledDate)
    .map(task => generateTaskEvent(task));
  
  return [ICS_HEADER, ...events, ICS_FOOTER].join('\r\n');
};

/**
 * 下载ICS文件
 * @param content ICS文件内容
 * @param filename 文件名
 */
export const downloadICSFile = (content: string, filename: string = 'todoing-tasks.ics') => {
  const blob = new Blob([content], { type: 'text/calendar;charset=utf-8' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
};