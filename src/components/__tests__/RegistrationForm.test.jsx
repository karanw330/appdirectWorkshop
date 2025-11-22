import { render, screen, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import RegistrationForm from '../RegistrationForm'
import * as api from '../../services/api'

// Mock the API
vi.mock('../../services/api', () => ({
  attendeesAPI: {
    getCount: vi.fn(),
    register: vi.fn(),
  },
}))

describe('RegistrationForm', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    api.attendeesAPI.getCount.mockResolvedValue({ data: { count: 0 } })
  })

  it('renders registration form', () => {
    render(<RegistrationForm />)
    expect(screen.getByText(/Register Now/i)).toBeInTheDocument()
    expect(screen.getByPlaceholderText(/Enter your full name/i)).toBeInTheDocument()
    expect(screen.getByPlaceholderText(/Enter your email/i)).toBeInTheDocument()
  })

  it('displays attendee count', async () => {
    api.attendeesAPI.getCount.mockResolvedValue({ data: { count: 42 } })
    render(<RegistrationForm />)
    
    await waitFor(() => {
      expect(screen.getByText('42')).toBeInTheDocument()
    })
  })

  it('validates required fields on submit', async () => {
    render(<RegistrationForm />)
    const submitButton = screen.getByRole('button', { name: /Register/i })
    
    submitButton.click()
    
    await waitFor(() => {
      expect(api.attendeesAPI.register).not.toHaveBeenCalled()
    })
  })

  it('submits form with valid data', async () => {
    api.attendeesAPI.register.mockResolvedValue({ data: { id: '123', name: 'Test User' } })
    
    render(<RegistrationForm />)
    
    const nameInput = screen.getByPlaceholderText(/Enter your full name/i)
    const emailInput = screen.getByPlaceholderText(/Enter your email/i)
    const designationSelect = screen.getByRole('combobox')
    const submitButton = screen.getByRole('button', { name: /Register/i })
    
    nameInput.value = 'Test User'
    emailInput.value = 'test@example.com'
    designationSelect.value = 'Software Engineer'
    
    submitButton.click()
    
    await waitFor(() => {
      expect(api.attendeesAPI.register).toHaveBeenCalledWith({
        name: 'Test User',
        email: 'test@example.com',
        designation: 'Software Engineer',
      })
    })
  })
})

