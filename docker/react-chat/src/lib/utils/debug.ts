// Debug utility to control console logging with rate limiting and categories
interface DebugConfig {
  enabled: boolean;
  categories: {
    connection: boolean;
    streaming: boolean;
    general: boolean;
    error: boolean;
    warn: boolean;
    critical: boolean; // For critical events that should always show
  };
  rateLimiting: {
    enabled: boolean;
    maxLogsPerSecond: number;
    windowSizeMs: number;
  };
}

class DebugManager {
  private config: DebugConfig = {
    enabled: true,
    categories: {
      connection: true,
      streaming: false,  // Reduce streaming noise
      general: true,
      error: true,
      warn: true,
      critical: true
    },
    rateLimiting: {
      enabled: true,
      maxLogsPerSecond: 10,
      windowSizeMs: 1000
    }
  };

  private logCounts = new Map<string, number>();
  private lastReset = Date.now();
  // Special tracking for critical messages to prevent spam
  private criticalMessageTracker = new Map<string, { count: number; lastLogged: number; firstLogged: number }>();

  private shouldLog(category: keyof DebugConfig['categories'], message: string): boolean {
    if (!this.config.enabled || !this.config.categories[category]) {
      return false;
    }

    if (!this.config.rateLimiting.enabled) {
      return true;
    }

    const now = Date.now();
    const key = `${category}:${message.substring(0, 50)}`; // Use message prefix as key

    // Reset counters if window has expired
    if (now - this.lastReset > this.config.rateLimiting.windowSizeMs) {
      this.logCounts.clear();
      this.lastReset = now;
    }

    // Special handling for critical messages to reduce spam
    if (category === 'critical') {
      return this.shouldLogCriticalMessage(message, now);
    }

    const currentCount = this.logCounts.get(key) || 0;
    if (currentCount >= this.config.rateLimiting.maxLogsPerSecond) {
      return false; // Rate limited
    }

    this.logCounts.set(key, currentCount + 1);
    return true;
  }

  private shouldLogCriticalMessage(message: string, now: number): boolean {
    // For critical messages, we want to be more sophisticated about rate limiting
    // to prevent spam while still showing important information
    
    // Create a more specific key for critical messages
    const key = message.substring(0, 100); // Use more of the message for better differentiation
    const tracker = this.criticalMessageTracker.get(key);
    
    // If this is a new message, log it
    if (!tracker) {
      this.criticalMessageTracker.set(key, { count: 1, lastLogged: now, firstLogged: now });
      return true;
    }
    
    // If we've logged this message recently, skip it
    const timeSinceLastLog = now - tracker.lastLogged;
    const timeSinceFirstLog = now - tracker.firstLogged;
    
    // If it's been less than 5 seconds since we last logged this message, skip it
    if (timeSinceLastLog < 5000) {
      // But still update the count
      this.criticalMessageTracker.set(key, { 
        count: tracker.count + 1, 
        lastLogged: tracker.lastLogged,
        firstLogged: tracker.firstLogged
      });
      return false;
    }
    
    // If it's been a while, log it but with a count of how many times we've seen it
    const newCount = tracker.count + 1;
    this.criticalMessageTracker.set(key, { 
      count: newCount, 
      lastLogged: now,
      firstLogged: tracker.firstLogged
    });
    
    // If this is a repeated message, modify the log to indicate how many times it's been repeated
    // and how long it's been repeating
    if (newCount > 1) {
      const duration = Math.floor(timeSinceFirstLog / 1000); // in seconds
      console.log(`[CRITICAL] ${message} (repeated ${newCount} times over ${duration}s)`);
      return false; // We've already logged it with the count, so don't log the original
    }
    
    return true;
  }

  // Public methods to control debug settings
  enableCategory(category: keyof DebugConfig['categories'], enabled: boolean = true) {
    this.config.categories[category] = enabled;
    console.log(`Debug category '${category}' ${enabled ? 'enabled' : 'disabled'}`);
  }

  enableDebug(enabled: boolean = true) {
    this.config.enabled = enabled;
    console.log(`Debug logging ${enabled ? 'enabled' : 'disabled'}`);
  }

  enableRateLimit(enabled: boolean = true, maxLogsPerSecond: number = 10) {
    this.config.rateLimiting.enabled = enabled;
    this.config.rateLimiting.maxLogsPerSecond = maxLogsPerSecond;
    console.log(`Rate limiting ${enabled ? 'enabled' : 'disabled'} (${maxLogsPerSecond}/sec)`);
  }

  getConfig() {
    return { ...this.config };
  }

  log(message: string, ...args: unknown[]) {
    if (this.shouldLog('general', message)) {
      console.log(`[DEBUG] ${message}`, ...args);
    }
  }

  error(message: string, ...args: unknown[]) {
    if (this.shouldLog('error', message)) {
      console.error(`[ERROR] ${message}`, ...args);
    }
  }

  warn(message: string, ...args: unknown[]) {
    if (this.shouldLog('warn', message)) {
      console.warn(`[WARN] ${message}`, ...args);
    }
  }

  connection(message: string, ...args: unknown[]) {
    if (this.shouldLog('connection', message)) {
      console.log(`[CONNECTION] ${message}`, ...args);
    }
  }

  streaming(message: string, ...args: unknown[]) {
    if (this.shouldLog('streaming', message)) {
      console.log(`[STREAMING] ${message}`, ...args);
    }
  }

  critical(message: string, ...args: unknown[]) {
    if (this.shouldLog('critical', message)) {
      // Special handling for critical messages that have been repeated
      const key = message.substring(0, 100);
      const tracker = this.criticalMessageTracker.get(key);
      
      if (tracker && tracker.count > 1) {
        // Already handled in shouldLogCriticalMessage, don't log again
        return;
      }
      
      console.log(`[CRITICAL] ${message}`, ...args);
    }
  }
}

// Create singleton instance
const debugManager = new DebugManager();

// Export the debug interface
export const debug = debugManager;

// Make debug controls available globally for development
if (typeof window !== 'undefined' && process.env.NODE_ENV === 'development') {
  (window as any).debug = {
    enable: (enabled?: boolean) => debugManager.enableDebug(enabled),
    enableCategory: (category: string, enabled?: boolean) => {
      const validCategories = ['connection', 'streaming', 'general', 'error', 'warn', 'critical'];
      if (validCategories.includes(category)) {
        debugManager.enableCategory(category as keyof DebugConfig['categories'], enabled);
      } else {
        console.warn(`Invalid debug category: ${category}`);
      }
    },
    enableRateLimit: (enabled?: boolean, maxLogs?: number) => debugManager.enableRateLimit(enabled, maxLogs),
    getConfig: () => debugManager.getConfig(),
    // Quick shortcuts
    enableConnection: () => debugManager.enableCategory('connection', true),
    disableConnection: () => debugManager.enableCategory('connection', false),
    enableStreaming: () => debugManager.enableCategory('streaming', true),
    disableStreaming: () => debugManager.enableCategory('streaming', false)
  };
}