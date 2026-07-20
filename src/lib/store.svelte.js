class SharedStore {
	/** @type {any[]} */
	orders = $state([]);
	/** @type {any[]} */
	sales = $state([]);
	/** @type {any[]} */
	expenses = $state([]);

	async init() {
		await Promise.all([this.fetchOrders(), this.fetchSales(), this.fetchExpenses()]);
		await this.checkCancelDates();
	}

	async fetchOrders() {
		try {
			const res = await fetch('http://localhost:8080/api/orders');
			if (res.ok) {
				this.orders = await res.json();
			}
		} catch (err) {
			console.error('Failed to fetch orders:', err);
		}
	}

	async fetchSales() {
		try {
			const res = await fetch('http://localhost:8080/api/sales');
			if (res.ok) {
				this.sales = await res.json();
			}
		} catch (err) {
			console.error('Failed to fetch sales:', err);
		}
	}

	async fetchExpenses() {
		try {
			const res = await fetch('http://localhost:8080/api/expenses');
			if (res.ok) {
				this.expenses = await res.json();
			}
		} catch (err) {
			console.error('Failed to fetch expenses:', err);
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
				updates.push(this.updateOrderStatus(order.id, 'Expired'));
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
			const res = await fetch('http://localhost:8080/api/orders', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(order)
			});
			if (res.ok) {
				await this.fetchOrders();
				await this.checkCancelDates();
			}
		} catch (err) {
			console.error('Failed to add order:', err);
		}
	}

	/**
	 * @param {string|number} id
	 * @param {string} [userRole]
	 */
	async removeOrder(id, userRole = 'Employee') {
		try {
			const res = await fetch(`http://localhost:8080/api/orders/${id}`, {
				method: 'DELETE',
				headers: { 'X-User-Role': userRole }
			});
			if (res.ok) {
				await this.fetchOrders();
			} else {
				const errData = await res.json();
				alert(errData.error || 'Failed to delete order.');
			}
		} catch (err) {
			console.error('Failed to delete order:', err);
		}
	}

	/**
	 * @param {string|number} id
	 * @param {string} newStatus
	 * @param {string} [invoiceId]
	 * @param {string} [updatedBy]
	 * @param {string} [userRole]
	 */
	async updateOrderStatus(id, newStatus, invoiceId, updatedBy, userRole = 'Employee') {
		try {
			const res = await fetch(`http://localhost:8080/api/orders/${id}/status`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					'X-User-Role': userRole
				},
				body: JSON.stringify({ status: newStatus, invoiceId, updatedBy })
			});
			if (res.ok) {
				await Promise.all([this.fetchOrders(), this.fetchSales()]);
			} else {
				const errData = await res.json();
				alert(errData.error || 'Failed to update order status.');
			}
		} catch (err) {
			console.error('Failed to update order status:', err);
		}
	}

	/**
	 * @param {any} expense
	 */
	async addExpense(expense) {
		try {
			const res = await fetch('http://localhost:8080/api/expenses', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(expense)
			});
			if (res.ok) {
				await this.fetchExpenses();
			}
		} catch (err) {
			console.error('Failed to add expense:', err);
		}
	}

	/**
	 * @param {string|number} id
	 * @param {string} [userRole]
	 */
	async removeExpense(id, userRole = 'Employee') {
		try {
			const res = await fetch(`http://localhost:8080/api/expenses/${id}`, {
				method: 'DELETE',
				headers: { 'X-User-Role': userRole }
			});
			if (res.ok) {
				await this.fetchExpenses();
			} else {
				const errData = await res.json();
				alert(errData.error || 'Failed to delete expense.');
			}
		} catch (err) {
			console.error('Failed to delete expense:', err);
		}
	}
}

export const sharedStore = new SharedStore();
