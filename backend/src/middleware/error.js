const errorHandler = (err, req, res, next) => {
  console.error(err.stack);
  
  // Mongoose bad ObjectId
  if (err.name === 'CastError') {
    return res.status(404).json({ msg: 'Resource not found' });
  }
  
  // Mongoose duplicate key
  if (err.code === 11000) {
    return res.status(400).json({ msg: 'Duplicate field value entered' });
  }
  
  // Mongoose validation error
  if (err.name === 'ValidationError') {
    const errors = Object.values(err.errors).map(e => e.message);
    return res.status(400).json({
      msg: 'Validation failed',
      errors
    });
  }
  
  res.status(500).json({
    success: false,
    error: 'Server Error'
  });
};

module.exports = errorHandler;