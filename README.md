# foot_event
一款go语言开发的容器各类事件采集上报的组件，该组件作用如下：<p>
1）实时采集监控kubernetes、mesos+marathon等容器资源编排各类变化事件，
为应用编排和应用调度运行增强运维能力；<p>
**采集源：**<p>
a.kubernetes通过API Server提供的各类事件服务API；<p>
b.mesos资源池之上marathon提供的各类事件服务API；<p>
**输出目标：**<p>
a.数据库mysql、oracle等；<p>
b.消息中间件kafka等；<p>
c.搜索引擎存储ES等；<p>
d.支持可扩展的采集对接。


