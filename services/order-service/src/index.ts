import express, { Request, Response } from 'express';

const app = express();
app.use(express.json());

interface Order {
  id: string;
  userId: string;
  itemId: string;
  quantity: number;
  totalAmount: number;
  status: 'PENDING' | 'CONFIRMED' | 'FAILED';
  createdAt: string;
}

const orders: Order[] = [];

// Readiness & Liveness probes
app.get('/healthz', (req: Request, res: Response) => {
  res.status(200).json({ status: 'UP', service: 'order-service' });
});

app.get('/metrics', (req: Request, res: Response) => {
  res.set('Content-Type', 'text/plain');
  res.send(`# HELP order_count_total Total count of orders processed
# TYPE order_count_total counter
order_count_total ${orders.length}
`);
});

// Order Creation Endpoint
app.post('/orders', (req: Request, res: Response) => {
  const { userId, itemId, quantity, pricePerUnit } = req.body;

  if (!userId || !itemId || !quantity) {
    return res.status(400).json({ error: 'Missing required order fields' });
  }

  const newOrder: Order = {
    id: `ord_${Date.now()}_${Math.floor(Math.random() * 1000)}`,
    userId,
    itemId,
    quantity,
    totalAmount: (pricePerUnit || 10) * quantity,
    status: 'CONFIRMED',
    createdAt: new Date().toISOString()
  };

  orders.push(newOrder);

  return res.status(201).json({
    message: 'Order created successfully',
    order: newOrder
  });
});

app.get('/orders', (req: Request, res: Response) => {
  return res.json(orders);
});

const PORT = process.env.PORT || 8083;
app.listen(PORT, () => {
  console.log(`[OrderService] Running on port ${PORT}`);
});
