const fastify = require('fastify')({ logger: true })

fastify.post('/', async (request, reply) => {
  const { text } = request.body
  
  if (!text) {
    return reply.status(400).send({ error: 'text is required' })
  }

  try {
    const response = await fetch('http://localhost:8080', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ text })
    })

    if (!response.ok) {
      throw new Error('Main server error')
    }

    const data = await response.json()
    return reply.send(data)
  } catch (error) {
    request.log.error(error)
    return reply.status(500).send({ error: 'Internal server error' })
  }
})

const start = async () => {
  try {
    await fastify.listen({ port: 3000 })
    fastify.log.info('Gateway running on port 3000')
  } catch (err) {
    fastify.log.error(err)
    process.exit(1)
  }
}

start()