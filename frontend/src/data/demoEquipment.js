export const demoEquipment = [
  {
    id: 101,
    name: "Sony A7 III с объективом 24-70",
    description:
      "Полный комплект для съемки: камера, объектив, аккумуляторы, карта памяти и сумка. Подойдет для репортажей, предметной съемки и учебных проектов.",
    image_url:
      "https://images.unsplash.com/photo-1516035069371-29a1b244cc32?auto=format&fit=crop&w=1200&q=80",
    location: "Екатеринбург",
    price_per_day: 1800,
    available: true,
    category: { name: "Фото и видео" },
    owner: {
      id: 7,
      name: "Алексей",
      city: "Екатеринбург",
      avatar_url: "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=300&q=80",
      bio: "Сдаю технику для учебных и коммерческих съемок."
    }
  },
  {
    id: 102,
    name: "Проектор Epson Full HD",
    description:
      "Яркий проектор для презентаций, лекций и домашних кинопоказов. В комплекте HDMI-кабель и пульт.",
    image_url:
      "https://images.unsplash.com/photo-1601944179066-29786cb9d32a?auto=format&fit=crop&w=1200&q=80",
    location: "Екатеринбург",
    price_per_day: 950,
    available: true,
    category: { name: "Презентации" },
    owner: {
      id: 8,
      name: "Мария",
      city: "Екатеринбург",
      avatar_url: "https://images.unsplash.com/photo-1494790108377-be9c29b29330?auto=format&fit=crop&w=300&q=80",
      bio: "Помогаю с оборудованием для мероприятий."
    }
  },
  {
    id: 103,
    name: "DJI Ronin RS 3",
    description:
      "Стабилизатор для плавной видеосъемки. Настроен, заряжен, есть быстрый инструктаж перед передачей.",
    image_url:
      "https://images.unsplash.com/photo-1616243850903-7b51366c1ba2?auto=format&fit=crop&w=1200&q=80",
    location: "Пермь",
    price_per_day: 1200,
    available: false,
    category: { name: "Фото и видео" },
    owner: {
      id: 9,
      name: "Илья",
      city: "Пермь",
      avatar_url: "https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?auto=format&fit=crop&w=300&q=80",
      bio: "Видеограф, сдаю свободное оборудование."
    }
  }
];
