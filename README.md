# ⚔️ ABYSS HUNTER - Dungeon Crawler PWA

## 🎮 Game Overview

**Abyss Hunter** adalah semi-idle roguelite dungeon crawler berbasis PWA dengan Dark Chibi Fantasy aesthetic. Game ini dirancang untuk memberikan dopamine progression yang addictive dengan sistem loot, refine, dan prestige.

### 📱 Platform Support
- ✅ Mobile Browser (iOS & Android)
- ✅ Desktop Browser (Chrome, Firefox, Edge)
- ✅ Installable PWA
- ✅ Partial Offline Play (Combat, Inventory Management)

---

## 🚀 Quick Start

### Prerequisites
- Node.js 16+ dan npm 8+
- Modern web browser dengan ES6 support

### Installation

```bash
# Clone atau download project
cd abyss-hunter

# Install dependencies
npm install

# Development mode
npm run dev

# Production build
npm run build

# Preview build
npm run preview
```

Server akan berjalan di `http://localhost:3000`

---

## 📁 Project Structure

```
abyss-hunter/
├── public/
│   ├── index.html          # HTML entry point
│   ├── sw.js               # Service Worker
│   └── manifest.json       # PWA manifest
├── src/
│   ├── App.jsx            # Root component
│   ├── App.css            # Global styles
│   ├── main.jsx           # React entry
│   ├── pages/
│   │   ├── GameMain.jsx   # Main game page
│   │   └── GameMain.css
│   ├── components/
│   │   ├── ClassSelection.jsx
│   │   ├── Battle.jsx
│   │   ├── Inventory.jsx
│   │   ├── Shop.jsx
│   │   ├── Blacksmith.jsx
│   │   ├── DungeonSelect.jsx
│   │   └── [Component].css
│   └── store/
│       ├── store.js       # Redux store config
│       └── slices/
│           ├── playerSlice.js
│           ├── inventorySlice.js
│           ├── battleSlice.js
│           ├── dungeonSlice.js
│           ├── shopSlice.js
│           ├── blacksmithSlice.js
│           └── prestigeSlice.js
├── package.json
├── vite.config.js
└── README.md
```

---

## 🎯 Gameplay Systems

### 1. **Class Selection**
5 kelas dengan specialty berbeda:
- **Knight** 🛡️ - Tank/Sustain
- **Assassin** 🗡️ - Crit/Dodge
- **Mage** 🔮 - Burst/AoE
- **Archer** 🏹 - Rapid Attack
- **Berserker** ⚡ - High Damage/Risk

### 2. **Combat System**
- **Auto-battle** realtime dengan manual attack option
- Damage calculation berdasarkan stats (STR, AGI, INT, VIT, LUCK)
- Crit system dengan adjustable speed
- Battle log untuk tracking damage

### 3. **Loot System**
- Random rarity drops (Common, Rare, Epic, Legendary, Mythic)
- Loot dopamine dengan visual feedback
- Drop rate bonus dari stats LUCK
- Item auto-collected ke inventory

### 4. **Equipment System**
- 9 equipment slots (Weapon, Armor, etc)
- Randomized stats pada setiap item
- Equipment rarity progression
- Quick equip/unequip dari inventory

### 5. **Refine System**
- Refine level +1 hingga +15
- Cost dan success rate meningkat tiap level
- Downgrade risk di level tinggi (+11-15)
- Material: Gold dan Refine Stone

### 6. **Shop/Merchant**
- Random merchant slots yang refresh
- Buy items dengan Gold
- Refresh merchant dengan Fragment (✨)
- Harga refresh meningkat setiap refresh

### 7. **Dungeon Progression**
- 3 dungeon dengan difficulty berbeda
- Floor progression dengan enemy scaling
- Max floor tracking per dungeon
- Boss fight setiap 10 floor

### 8. **Save System**
- Auto-save dengan Redux-persist
- LocalStorage persistence
- Offline inventory/equipment management
- Session save di browser memory

---

## 💡 Core Game Loop

```
1. SELECT CLASS → 2. CHOOSE DUNGEON → 3. ENTER BATTLE
    ↓ (Loop)
4. AUTO COMBAT → 5. GET LOOT → 6. LEVEL UP
    ↓
7. NEXT FLOOR (atau DEATH/EXIT) → 8. INVENTORY MANAGEMENT
    ↓
9. SHOP → 10. BLACKSMITH (REFINE) → kembali ke step 1
```

---

## 🎮 Game Controls

### Battle Screen
- **Attack Button** - Manual attack (saat auto off)
- **Auto/Stop Button** - Toggle auto-battle
- **Speed Slider** - Adjust battle speed (0.5x - 2.5x)
- **Next Floor** - Lanjut ke floor berikutnya
- **Exit Run** - Keluar dungeon

### Navigation
- **Tab Buttons** - Switch antar fitur
- **Equipment Slots** - Click untuk equip/unequip
- **Sort Buttons** - Sorting inventory (Recent, Rarity, Name)

---

## 📊 Progression System

### Player Stats
```
STR (Strength)     → Physical Damage
AGI (Agility)      → Attack Speed + Dodge
INT (Intelligence) → Magic Damage
VIT (Vitality)     → HP + Defense
LUCK (Luck)        → Crit Chance + Drop Rate
```

### Level Up Progression
- Exp untuk level up meningkat exponential (1.1x multiplier)
- Level up memberikan +1-2 per stat
- Max HP meningkat dengan VIT scaling
- Tidak ada level cap

### Prestige (Phase 2)
- Reset semua tapi dapat Soul Point
- Unlock Soul Tree untuk permanent buffs
- EXP Bonus, Drop Rate, Damage, Merchant Luck

---

## ⚙️ Tech Stack

### Frontend
- **React 18** - UI Framework
- **Redux Toolkit** - State Management
- **Redux-Persist** - Local Storage Persistence
- **Vite** - Build tool
- **Tailwind CSS** - Styling (via custom CSS)

### PWA
- **Service Worker** - Offline capability
- **Web App Manifest** - Installation
- **LocalStorage** - Data persistence

### Performance
- < 150MB RAM browser
- 60fps target on mobile
- Fast loading dengan code splitting

---

## 🔧 Development Guide

### Adding New Features

#### 1. New Dungeon
```javascript
// dungeonSlice.js
dungeons: {
  new_dungeon: {
    name: 'Name',
    theme: 'Theme',
    maxFloors: 100,
    currentFloor: 1,
    maxReachedFloor: 1,
  }
}
```

#### 2. New Item/Equipment
```javascript
// Modify SHOP_ITEMS atau generateLoot
const newItem = {
  name: 'Item Name',
  rarity: 'rare',
  type: 'equipment',
  stats: {
    attack: 10,
    defense: 5,
  }
}
```

#### 3. New Skill/Ability
```javascript
// Create new battle logic
const skill = {
  name: 'Skill Name',
  damage: baseDamage * multiplier,
  cooldown: 5000,
  effect: 'effect description'
}
```

### Modifying Game Balance

```javascript
// battle.js - Enemy scaling
const baseHp = 20 + floor * 5;        // Adjust multiplier
const baseDamage = 3 + floor * 1.5;   // Adjust scaling

// Battle.jsx - Loot chance
if (roll < 0.6) rarity = 'common';    // Adjust drop rates
else if (roll < 0.85) rarity = 'rare';
```

---

## 📱 PWA Installation

### Desktop (Chrome, Edge)
1. Buka game di browser
2. Click "Install" di address bar
3. Atau klik menu → "Install app"

### Mobile iOS
1. Buka di Safari
2. Click Share → "Add to Home Screen"
3. Beri nama dan install

### Mobile Android
1. Buka di Chrome
2. Click menu → "Install app"
3. Atau tap install prompt di bottom

---

## 🐛 Debugging

### Redux DevTools
```javascript
// Akan terintegrasi dengan Redux DevTools extension jika installed
// Lihat state di: chrome://extensions → Redux DevTools
```

### LocalStorage Check
```javascript
// Browser Console
localStorage.getItem('persist:root')  // Cek stored state
Object.keys(localStorage)             // List semua storage keys
```

### Service Worker Debug
```javascript
// Chrome DevTools → Application → Service Workers
// Lihat registered SW dan cache storage
```

---

## 📈 Performance Optimization

### Already Implemented
- ✅ Lazy loading components
- ✅ Memoized selectors
- ✅ CSS-only animations
- ✅ Efficient re-renders
- ✅ Service Worker caching

### Future Optimization
- Code splitting per route
- Image optimization
- WebP format support
- Progressive JPEG loading
- Bundle analysis

---

## 🎨 Customization

### Colors (CSS Variables)
Edit di `App.css`:
```css
:root {
  --color-primary: #00d4ff;        /* Cyan */
  --color-accent: #ff006e;         /* Pink */
  --color-accent-alt: #ffd60a;     /* Yellow */
  /* ... more colors */
}
```

### Fonts
```css
@import url('https://fonts.googleapis.com/css2?family=Fredoka:wght@400;600;700&family=Pixelify+Sans:wght@400;600;700&display=swap');
```

### Animations
Semua animations di `App.css` dengan CSS keyframes

---

## 📝 MVP Features (Phase 1)

✅ **Completed:**
- 2+ class options (5 kelas)
- 1 main dungeon (3 dungeons)
- 30+ floor progression
- Auto combat system
- Loot mechanics
- Inventory system
- Equipment system
- Refine system
- Merchant shop
- Blacksmith
- Save/Load system
- Responsive design
- PWA functionality

### Phase 2 (Planned)
- Rune system
- Relic system
- Prestige system
- Element system
- Set items
- Advanced events

### Phase 3 (Planned)
- Endless tower
- Daily modifier
- Advanced AI
- Seasonal events
- Leaderboard

---

## 🚀 Deployment

### Deploy ke Vercel (Recommended)
```bash
npm i -g vercel
vercel
```

### Deploy ke Netlify
```bash
# Build dulu
npm run build

# Upload dist folder ke Netlify
# Atau connect GitHub repo
```

### Deploy ke GitHub Pages
```bash
npm run build
# Upload dist ke gh-pages branch
```

---

## 🐛 Known Issues & Limitations

1. **Offline Mode** - Combat hanya semi-offline, perlu koneksi untuk beberapa fitur
2. **Mobile Zoom** - Zoom disabled untuk prevent UI breaking
3. **Screen Orientation** - Best at portrait untuk mobile
4. **Storage Limit** - LocalStorage ~10MB limit

---

## 💬 Game Tips for Players

1. **Early Game** - Focus farming gold untuk equipment
2. **Mid Game** - Refine equipment untuk stat boost
3. **Loot Priority** - Rare+ item lebih valuable
4. **Floor Progression** - Setiap floor 10x lebih sulit
5. **Merchant Luck** - Sering refresh shop untuk item bagus
6. **Equipment Balance** - Mix stat untuk optimal build

---

## 📚 Resources & References

### Redux Documentation
- https://redux-toolkit.js.org/
- https://react-redux.js.org/

### PWA Best Practices
- https://web.dev/progressive-web-apps/
- https://developer.mozilla.org/en-US/docs/Web/Progressive_web_apps

### React Docs
- https://react.dev

### Vite Guide
- https://vitejs.dev/

---

## 📄 License

MIT License - Feel free to use, modify, dan distribute

---

## 🎉 Credits

Created with ❤️ for the roguelite gaming community.

**Abyss Hunter** - Where every run matters, and dopamine never stops.

---

## 📞 Support & Feedback

Found a bug? Have a feature request?
- Open an issue
- Submit feedback
- Join the community

**Happy hunting!** ⚔️✨
