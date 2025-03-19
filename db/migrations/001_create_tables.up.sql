-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    phone VARCHAR(20),
    avatar VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 创建用户偏好表
CREATE TABLE IF NOT EXISTS user_preferences (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    language VARCHAR(10) NOT NULL DEFAULT 'zh-CN',
    currency VARCHAR(3) NOT NULL DEFAULT 'CNY',
    distance_unit VARCHAR(10) NOT NULL DEFAULT 'km',
    temperature_unit VARCHAR(10) NOT NULL DEFAULT 'celsius',
    notification_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    theme VARCHAR(10) NOT NULL DEFAULT 'light'
);

-- 创建景点表
CREATE TABLE IF NOT EXISTS attractions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    city VARCHAR(50) NOT NULL,
    country VARCHAR(50) NOT NULL,
    address TEXT,
    latitude DECIMAL(10, 7) NOT NULL,
    longitude DECIMAL(10, 7) NOT NULL,
    image_url TEXT,
    category VARCHAR(50) NOT NULL,  -- 博物馆、自然景观、历史遗迹等
    tags JSONB,  -- 存储标签数组，如["文化", "艺术", "古迹"]
    open_hours TEXT,  -- 存储营业时间信息
    ticket_price DECIMAL(10, 2),  -- 门票价格
    duration INTEGER,  -- 推荐游览时长（分钟）
    rating DECIMAL(3, 1),  -- 评分，1-5
    popularity INTEGER DEFAULT 0,  -- 人气指数
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_attractions_city ON attractions(city);
CREATE INDEX idx_attractions_country ON attractions(country);
CREATE INDEX idx_attractions_category ON attractions(category);
CREATE INDEX idx_attractions_rating ON attractions(rating);

-- 创建行程表
CREATE TABLE IF NOT EXISTS itineraries (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    destination VARCHAR(100) NOT NULL,  -- 目的地/城市
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,  -- 是否公开分享
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 创建行程天表
CREATE TABLE IF NOT EXISTS itinerary_days (
    id SERIAL PRIMARY KEY,
    itinerary_id INTEGER REFERENCES itineraries(id) ON DELETE CASCADE,
    day_number INTEGER NOT NULL,  -- 第几天
    date DATE NOT NULL,
    note TEXT
);

-- 创建行程项目表（景点、酒店、餐厅等）
CREATE TABLE IF NOT EXISTS itinerary_items (
    id SERIAL PRIMARY KEY,
    itinerary_day_id INTEGER REFERENCES itinerary_days(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL,  -- 类型：attraction, hotel, transport, meal, etc.
    ref_id INTEGER,  -- 引用ID，如景点ID
    title VARCHAR(100) NOT NULL,
    description TEXT,
    start_time INTEGER,  -- 存储分钟数，如9:00存储为540
    end_time INTEGER,  -- 存储分钟数
    duration INTEGER,  -- 持续时间（分钟）
    item_order INTEGER NOT NULL,  -- 项目排序
    location TEXT,
    latitude DECIMAL(10, 7),
    longitude DECIMAL(10, 7),
    notes TEXT
);

-- 创建用户评价表
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    attraction_id INTEGER REFERENCES attractions(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, attraction_id)  -- 每个用户对每个景点只能评价一次
);

-- 创建索引
CREATE INDEX idx_itineraries_user_id ON itineraries(user_id);
CREATE INDEX idx_itineraries_is_public ON itineraries(is_public);
CREATE INDEX idx_itinerary_days_itinerary_id ON itinerary_days(itinerary_id);
CREATE INDEX idx_itinerary_items_itinerary_day_id ON itinerary_items(itinerary_day_id);
CREATE INDEX idx_reviews_attraction_id ON reviews(attraction_id);
CREATE INDEX idx_reviews_user_id ON reviews(user_id); 