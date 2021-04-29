using System;
using StatsdClient;

namespace HelloApi.Loggers
{
    public class DataDogStatsdEventLogger : IEventLogger
    {
        private DogStatsdService _dataDogStatsdClient;
        private bool _disposedValue;

        public DataDogStatsdEventLogger(StatsdConfig config)
        {
            _dataDogStatsdClient = new DogStatsdService();
            _dataDogStatsdClient.Configure(config);
        }

        public void LogEvent(string eventName, string eventData)
        {
            _dataDogStatsdClient.Event(eventName, eventData.ToString());
        }

        protected virtual void Dispose(bool disposing)
        {
            if (!_disposedValue)
            {
                if (disposing)
                {
                    _dataDogStatsdClient?.Dispose();
                }

                _disposedValue = true;
            }
        }

        public void Dispose()
        {
            Dispose(disposing: true);
            GC.SuppressFinalize(this);
        }
    }
}
