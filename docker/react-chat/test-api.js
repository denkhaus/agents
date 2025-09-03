// Test script to verify API connection
async function testAPI() {
  try {
    console.log('Testing API connection...')
    
    const response = await fetch('http://localhost:6999/list-apps', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    })
    
    console.log('Response status:', response.status)
    console.log('Response headers:', [...response.headers.entries()])
    
    if (response.ok) {
      const data = await response.json()
      console.log('API data:', data)
      return data
    } else {
      console.error('API error:', response.status, response.statusText)
    }
  } catch (error) {
    console.error('Network error:', error)
  }
}

testAPI()