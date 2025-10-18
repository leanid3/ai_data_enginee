import React, { useState } from 'react';
import { PipelineProvider } from './contexts/PipelineContext';
import ErrorBoundary from './components/common/ErrorBoundary';
import Notification from './components/common/Notification';
import PipelineWizard from './components/wizard/PipelineWizard';
import PipelineList from './components/pipeline/PipelineList';
import './App.module.css';

function App() {
  const [currentView, setCurrentView] = useState('wizard');
  const [notifications, setNotifications] = useState([]);

  const showNotification = (message, type = 'info') => {
    const id = Date.now();
    const notification = { id, message, type };
    setNotifications(prev => [...prev, notification]);

    setTimeout(() => {
      setNotifications(prev => prev.filter(n => n.id !== id));
    }, 5000);
  };

  const removeNotification = (id) => {
    setNotifications(prev => prev.filter(n => n.id !== id));
  };

  const containerStyle = {
    minHeight: '100vh',
    backgroundColor: '#f5f6fa',
    fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
  };

  const headerStyle = {
    backgroundColor: '#2c3e50',
    color: 'white',
    padding: '20px 0',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
  };

  const headerContentStyle = {
    maxWidth: '1200px',
    margin: '0 auto',
    padding: '0 20px',
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
  };

  const titleStyle = {
    fontSize: '24px',
    fontWeight: 'bold',
    margin: 0,
  };

  const navStyle = {
    display: 'flex',
    gap: '16px',
  };

  const navButtonStyle = (isActive) => ({
    padding: '10px 20px',
    backgroundColor: isActive ? '#3498db' : 'transparent',
    color: 'white',
    border: '1px solid #3498db',
    borderRadius: '6px',
    cursor: 'pointer',
    fontSize: '14px',
    fontWeight: 'bold',
    transition: 'all 0.3s ease',
  });

  const mainStyle = {
    maxWidth: '1200px',
    margin: '0 auto',
    padding: '20px',
  };

  const footerStyle = {
    backgroundColor: '#34495e',
    color: 'white',
    padding: '20px 0',
    textAlign: 'center',
    marginTop: '40px',
  };

  const footerContentStyle = {
    maxWidth: '1200px',
    margin: '0 auto',
    padding: '0 20px',
  };

  return (
    <ErrorBoundary>
      <PipelineProvider>
        <div style={containerStyle}>
          <header style={headerStyle}>
            <div style={headerContentStyle}>
              <h1 style={titleStyle}>
                🚀 Интеллектуальный цифровой инженер данных
              </h1>
              <nav style={navStyle}>
          <button
                  style={navButtonStyle(currentView === 'wizard')}
                  onClick={() => setCurrentView('wizard')}
                >
                  📊 Создать пайплайн
              </button>
              <button
                  style={navButtonStyle(currentView === 'pipelines')}
                  onClick={() => setCurrentView('pipelines')}
              >
                  📋 Управление
              </button>
              </nav>
            </div>
          </header>

          <main style={mainStyle}>
            {currentView === 'wizard' && <PipelineWizard />}
            {currentView === 'pipelines' && <PipelineList />}
          </main>

          <footer style={footerStyle}>
            <div style={footerContentStyle}>
              <p style={{ margin: 0, fontSize: '14px' }}>
                © 2024 AI Data Engineer Backend. Все права защищены.
              </p>
            </div>
          </footer>

          {/* Уведомления */}
          {notifications.map(notification => (
            <Notification
              key={notification.id}
              message={notification.message}
              type={notification.type}
              onClose={() => removeNotification(notification.id)}
            />
          ))}
          </div>
      </PipelineProvider>
    </ErrorBoundary>
  );
}

export default App;