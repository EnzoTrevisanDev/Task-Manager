import React from 'react';

const UserSelect = ({ users, value, onChange, required = false, placeholder = 'Select a user', multiple = false, onRemoveUser }) => {
  console.log('UserSelect received users:', users);

  if (!users || users.length === 0) {
    return <div style={{ color: '#9ca3af', fontSize: '14px' }}>No users available</div>;
  }

  // Filter out any invalid users and ensure we have the correct data structure
  const validUsers = users.filter(user => {
    const hasValidId = user.id || user.UserID || user.ID;
    const hasValidName = user.name || user.Name || user.User?.name;
    return hasValidId && hasValidName;
  });

  console.log('Valid users:', validUsers);

  if (validUsers.length === 0) {
    return <div style={{ color: '#9ca3af', fontSize: '14px' }}>No valid users available</div>;
  }

  return (
    <div className="user-select">
      <div className="selected-users">
        {validUsers.map((user) => (
          <div key={`selected-${user.id}`} className="selected-user">
            <span>{user.name}</span>
            {onRemoveUser && (
              <button onClick={() => onRemoveUser(user.id)}>×</button>
            )}
          </div>
        ))}
      </div>
      <div className="user-dropdown">
        <select
          value={value}
          onChange={onChange}
          required={required}
          multiple={multiple}
          className="form-select"
          style={{ width: '100%', padding: '8px', borderRadius: '4px', border: '1px solid #d1d5db' }}
        >
          <option value="">{placeholder}</option>
          {validUsers.map((user) => {
            const userId = user.id || user.UserID || user.ID;
            const userName = user.name || user.Name || user.User?.name;
            console.log('Processing user in select:', { userId, userName });
            return (
              <option key={`user-${userId}`} value={userId}>
                {userName}
              </option>
            );
          })}
        </select>
      </div>
    </div>
  );
};

export default UserSelect; 