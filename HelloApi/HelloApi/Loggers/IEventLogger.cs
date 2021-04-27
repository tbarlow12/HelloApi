using System;
namespace HelloApi.Loggers
{
    public interface IEventLogger : IDisposable
    {
        void LogEvent(string eventName, string eventData);
    }
}
