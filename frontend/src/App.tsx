import { useEffect, useState } from 'react'
import './App.css'

// 1. 定義 Go 後端傳來的商品資料長什麼樣子 (TypeScript 的強項！)
interface Item {
  ID: number;
  name: string;
  price: number;
  status: string;
  image_url: string;
}

function App() {
  // 2. 準備兩個「狀態 (State)」來裝資料和畫面載入狀態
  const [items, setItems] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);

  // 3. useEffect：當這個網頁一打開的時候，要做什麼事？
  useEffect(() => {
    // 寫信去問你的 Go 伺服器拿資料！(記得你剛剛開的 CORS 通行證嗎？現在派上用場了)
    fetch('http://localhost:8080/items')
      .then((res) => res.json())
      .then((data) => {
        // Go 的 API 回傳的是 { source: "...", limit: 10, items: [...] }
        // 我們只要裡面的 items 陣列！
        setItems(data.items || []);
        setLoading(false);
      })
      .catch((err) => {
        console.error("無法取得資料：", err);
        setLoading(false);
      });
  }, []);

  // 4. 如果資料還在拿，顯示載入中
  if (loading) {
    return <h2>🔄 正在從 Go 伺服器極速拉取商品中...</h2>;
  }

  // 5. 資料拿到了！把商品變成一個個漂亮的卡片印在畫面上
  return (
    <div style={{ padding: '20px', fontFamily: 'sans-serif' }}>
      <h1>🛒 歡迎來到我的 Mercari Clone</h1>
      
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))', gap: '20px' }}>
        {items.map((item) => (
          <div key={item.ID} style={{ border: '1px solid #ccc', borderRadius: '8px', padding: '16px', textAlign: 'left' }}>
            {/* 注意：後端給的圖檔網址是 /uploads/xxx，所以要補上後端的 Port 8080 */}
            {item.image_url ? (
               <img 
                 src={`http://localhost:8080${item.image_url}`} 
                 alt={item.name} 
                 style={{ width: '100%', height: '200px', objectFit: 'cover', borderRadius: '4px' }} 
               />
            ) : (
               <div style={{ width: '100%', height: '200px', backgroundColor: '#eee', display: 'flex', alignItems: 'center', justifyContent: 'center', borderRadius: '4px' }}>無圖片</div>
            )}
            
            <h3>{item.name}</h3>
            <p style={{ color: '#e53e3e', fontWeight: 'bold', fontSize: '1.2rem' }}>¥ {item.price}</p>
            <p style={{ color: item.status === 'sold' ? 'red' : 'green' }}>
              狀態: {item.status === 'sold' ? '已售出' : '販售中'}
            </p>
          </div>
        ))}
      </div>
    </div>
  )
}

export default App