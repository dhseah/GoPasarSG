package mysql

// ShoppingCartCleanup cleans up all the shopping cart
// by removing items that were added more than 3 days ago.
// ShoppingCartCleanup is called by a gorountine which is
// started from main when the application is started.
func (m *CartModel) ShoppingCartCleanUp() error {
	_, err := m.DB.Exec(`DELETE FROM ShoppingCart WHERE Modified < (NOW() - INTERVAL 3 DAY)`)
	if err != nil {
		return err
	}
	return nil

}

// VerifiedUserCleanUp removes all unverified users who
// have not verified themselves within 7 days from signup.
// VerifiedUserCleanup is called by a gorountine which is
// started from main when the application is started
func (m *UserModel) VerifiedUserCleanUp() error {
	_, err := m.DB.Exec(`DELETE FROM User WHERE Created < (NOW() - INTERVAL 7 DAY) AND Verified = 0`)
	if err != nil {
		return err
	}

	return nil
}
