// Simple test script to verify API connectivity
const API_BASE_URL = 'http://localhost:6999';

async function testAPI() {
  console.log('Testing API connectivity...');
  
  try {
    // Test list apps endpoint
    console.log('Testing /list-apps endpoint...');
    const appsResponse = await fetch(`${API_BASE_URL}/list-apps`);
    console.log('Apps response status:', appsResponse.status);
    
    if (appsResponse.ok) {
      const apps = await appsResponse.json();
      console.log('Available apps:', apps);
      
      if (apps && apps.length > 0) {
        const agentId = apps[0];
        console.log(`Testing with agent: ${agentId}`);
        
        // Test create session
        console.log('Creating session...');
        const sessionResponse = await fetch(`${API_BASE_URL}/apps/${agentId}/users/testuser/sessions`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          }
        });
        
        console.log('Session response status:', sessionResponse.status);
        if (sessionResponse.ok) {
          const session = await sessionResponse.json();
          console.log('Created session:', session);
          
          // Test run agent
          console.log('Testing agent run...');
          const runResponse = await fetch(`${API_BASE_URL}/run_sse`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Accept': 'text/event-stream'
            },
            body: JSON.stringify({
              appName: agentId,
              userID: 'testuser',
              sessionID: session.id,
              streaming: true,
              newMessage: {
                role: 'user',
                parts: [{ text: 'Hello, test message' }]
              }
            })
          });
          
          console.log('Run response status:', runResponse.status);
          console.log('Run response headers:', Object.fromEntries(runResponse.headers.entries()));
          
          if (runResponse.ok && runResponse.body) {
            console.log('Streaming response received, reading first chunk...');
            const reader = runResponse.body.getReader();
            const decoder = new TextDecoder();
            
            try {
              const { done, value } = await reader.read();
              if (!done) {
                const chunk = decoder.decode(value, { stream: true });
                console.log('First chunk:', chunk.substring(0, 200));
              } else {
                console.log('Stream ended immediately');
              }
            } catch (readError) {
              console.error('Error reading stream:', readError);
            } finally {
              reader.releaseLock();
            }
          } else {
            console.log('Run failed or no body');
            const errorText = await runResponse.text().catch(() => 'Unable to read error response');
            console.log('Error response:', errorText);
          }
        } else {
          console.log('Session creation failed');
          const errorText = await sessionResponse.text().catch(() => 'Unable to read error response');
          console.log('Error response:', errorText);
        }
      }
    } else {
      console.log('Failed to get apps list');
      const errorText = await appsResponse.text().catch(() => 'Unable to read error response');
      console.log('Error response:', errorText);
    }
  } catch (error) {
    console.error('API test failed:', error);
  }
}

testAPI();