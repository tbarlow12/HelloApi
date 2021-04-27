using System;
namespace HelloApi.Loggers
{
    public interface IEventLogger : IDisposable
    {
        public void LogEvent(string eventName, string eventData);
    }
}
