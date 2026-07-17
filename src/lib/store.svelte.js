class SharedStore {
  /** @type {any[]} */
  orders = $state([]);
  /** @type {any[]} */
  sales = $state([]);

  async init() {
    await Promise.all([this.fetchOrders(), this.fetchSales()]);
    await this.checkCancelDates();
  }

  async fetchOrders() {
    try {
      const res = await fetch("http://localhost:8080/api/orders");
      if (res.ok) {
        this.orders = await res.json();
      }
    } catch (err) {
      console.error("Failed to fetch orders:", err);
    }
  }

  async fetchSales() {
    try {
      const res = await fetch("http://localhost:8080/api/sales");
      if (res.ok) {
        this.sales = await res.json();
      }
    } catch (err) {
      console.error("Failed to fetch sales:", err);
    }
  }

  /**
   * Automatically synchronizes and updates status of expired orders in database.
   */
  async checkCancelDates() {
    const todayStr = new Date().toISOString().split('T')[0];
    const updates = [];
    
    for (const order of this.orders) {
      if (
        order.cancelDate && 
        order.cancelDate !== 'N/A' && 
        order.cancelDate < todayStr && 
        order.status === 'Pending'
      ) {
        updates.push(this.updateOrderStatus(order.id, 'Cancelled'));
      }
    }
    
    if (updates.length > 0) {
      await Promise.all(updates);
    }
  }

  /**
   * @param {any} order
   */
  async addOrder(order) {
    try {
      const res = await fetch("http://localhost:8080/api/orders", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(order)
      });
      if (res.ok) {
        await this.fetchOrders();
        await this.checkCancelDates();
      }
    } catch (err) {
      console.error("Failed to add order:", err);
    }
  }

  /**
   * @param {string|number} id
   */
  async removeOrder(id) {
    try {
      const res = await fetch(`http://localhost:8080/api/orders/${id}`, {
        method: "DELETE"
      });
      if (res.ok) {
        await this.fetchOrders();
      }
    } catch (err) {
      console.error("Failed to delete order:", err);
    }
  }

  /**
   * @param {string|number} id
   * @param {string} newStatus
   * @param {string} [invoiceId]
   * @param {string} [updatedBy]
   */
  async updateOrderStatus(id, newStatus, invoiceId, updatedBy) {
    try {
      const res = await fetch(`http://localhost:8080/api/orders/${id}/status`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ status: newStatus, invoiceId, updatedBy })
      });
      if (res.ok) {
        await Promise.all([this.fetchOrders(), this.fetchSales()]);
      }
    } catch (err) {
      console.error("Failed to update order status:", err);
    }
  }
}

export const sharedStore = new SharedStore();
