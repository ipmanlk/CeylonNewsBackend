from flask import Flask, request, jsonify
import undetected_chromedriver as uc
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import time
import logging

app = Flask(__name__)

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class BrowserScraper:
    def __init__(self):
        self.driver = None
        
    def get_driver(self):
        if self.driver is None:
            options = uc.ChromeOptions()
            options.add_argument('--no-sandbox')
            options.add_argument('--disable-dev-shm-usage')
            options.add_argument('--disable-gpu')
            options.add_argument('--disable-images')
            options.add_argument('--disable-css')
            options.add_argument('--headless')
            options.add_argument('--window-size=1920,1080')
            options.add_argument('--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36')
            
            self.driver = uc.Chrome(options=options)
            self.driver.implicitly_wait(10)
        return self.driver
    
    def scrape(self, url, wait_time=3):
        driver = self.get_driver()
        
        try:
            logger.info(f"Scraping URL: {url}")
            driver.get(url)
            
            # Wait for page to load
            time.sleep(wait_time)
            
            # Block images and CSS after page loads
            driver.execute_script("""
                // Block images
                var images = document.getElementsByTagName('img');
                for(var i = 0; i < images.length; i++) {
                    images[i].style.display = 'none';
                }
                
                // Block CSS
                var links = document.getElementsByTagName('link');
                for(var i = 0; i < links.length; i++) {
                    if(links[i].rel === 'stylesheet') {
                        links[i].disabled = true;
                    }
                }
                
                // Block style tags
                var styles = document.getElementsByTagName('style');
                for(var i = 0; i < styles.length; i++) {
                    styles[i].disabled = true;
                }
            """)
            
            # Get the HTML
            html = driver.page_source
            
            return {
                "success": True,
                "html": html,
                "url": url
            }
            
        except Exception as e:
            logger.error(f"Error scraping {url}: {str(e)}")
            return {
                "success": False,
                "error": str(e),
                "url": url
            }

# Initialize scraper
scraper = BrowserScraper()

@app.route('/')
def home():
    return jsonify({
        "message": "Web Scraping API",
        "endpoint": "/scrape",
        "method": "GET",
        "params": {
            "url": "URL to scrape (required)",
            "wait_time": "Seconds to wait for page load (default: 3)"
        },
        "example": "/scrape?url=https://example.com&wait_time=5"
    })

@app.route('/scrape', methods=['GET'])
def scrape():
    try:
        url = request.args.get('url')
        wait_time = int(request.args.get('wait_time', 10))
        
        if not url:
            return jsonify({"error": "URL parameter is required"}), 400
        
        result = scraper.scrape(url, wait_time)
        
        if result['success']:
            return jsonify(result)
        else:
            return jsonify(result), 500
            
    except Exception as e:
        logger.error(f"API error: {str(e)}")
        return jsonify({"error": str(e)}), 500

@app.route('/health')
def health():
    return jsonify({
        "status": "healthy",
        "browser_ready": scraper.driver is not None
    })

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8000, debug=False)
